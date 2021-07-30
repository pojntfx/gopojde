package orchestration

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/alessio/shellescape"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
	"github.com/iancoleman/strcase"
	"github.com/pojntfx/gopojde/pkg/config"
	"github.com/pojntfx/gopojde/pkg/util"
)

const (
	containerPrefix            = "/pojde-"
	volumePrefix               = "pojde-"
	configurationScriptsDir    = "/opt/pojde/configuration"
	caCertFile                 = "ca.pem"
	pojdeDockerImage           = "pojntfx/pojde:latest"
	firstInternalPort          = 8000
	portCount                  = 6 // portSuffixToServiceMap should have this many entries
	startCmd                   = "/lib/systemd/systemd"
	execPerm                   = 0775
	preferencesDirInContainer  = "/opt/pojde/preferences"
	preferencesFileInContainer = "preferences.sh"
	sshKeysPath                = "/root/.ssh/authorized_keys"

	StatusPullingLatestImage        = "pullingLatestImage"
	StatusRemovingExistingContainer = "removingExistingContainer"
	StatusCreatingContainer         = "creatingContainer"
	StatusPreparingConfigScripts    = "preparingConfigScripts"
)

var (
	journalctlTailCmd                = []string{"journalctl", "-f"}
	enumerateConfigurationScriptsCmd = []string{"ls", configurationScriptsDir}
	bashCmd                          = []string{"bash"}
	blocklistedScripts               = []string{"parameters.sh"} // TODO: Remove this as soon as scripts are no longer interactive
	configScripts                    = []string{"user", "apt", "code-server", "ttyd", "novnc", "jupyter-lab", "nginx", "docker", "pojdectl", "ssh", "git", "modules", "init", "clean"}
	portToServiceMap                 = map[int32]string{
		firstInternalPort + 0: "cockpit",
		firstInternalPort + 1: "codeserver",
		firstInternalPort + 2: "ttyd",
		firstInternalPort + 3: "novnc",
		firstInternalPort + 4: "jupyterlab",
		firstInternalPort + 5: "ssh",
	}
)

func getStartingMessage(status string) string {
	return status + "Starting"
}

func getDoneMessage(status string) string {
	return status + "Done"
}

func GetConfigScriptStatus(configScriptName string) string {
	return "runningConfigScript" + strcase.ToCamel(configScriptName)
}

func removePrefix(instanceName string) string {
	return strings.TrimPrefix(instanceName, containerPrefix)
}

func addContainerPrefix(instanceName string) string {
	return containerPrefix + instanceName
}

func addVolumePrefix(instanceName string) string {
	return volumePrefix + instanceName
}

func getDEBCacheVolumeName(instanceName string) string {
	return addVolumePrefix(instanceName + "-apt-cache")
}

func getPreferencesVolumeName(instanceName string) string {
	return addVolumePrefix(instanceName + "-preferences")
}

func getUserDataVolumeNames(instanceName string) []string {
	return []string{addVolumePrefix(instanceName + "-home-root"), addVolumePrefix(instanceName + "-home-user")}
}

func getCAVolumeName() string {
	return addVolumePrefix("ca")
}

func getTransferDirectory(instanceName string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, "Documents", "pojde", instanceName), nil
}

func getCommandToExecuteRefreshScript(scriptPath string) []string {
	return []string{"bash", "-c", fmt.Sprintf(". %v && refresh", shellescape.Quote(scriptPath))}
}

func getCommandToExecuteUpgradeScript(scriptPath string) []string {
	return []string{"bash", "-c", fmt.Sprintf(". %v && upgrade", shellescape.Quote(scriptPath))}
}

type InstanceRemovalOptions struct {
	Customizations bool
	DEBCache       bool
	Preferences    bool
	UserData       bool
	Transfer       bool
}

type InstanceCreationFlags struct {
	StartPort       int32
	PullLatestImage bool
	Recreate        bool
	Isolate         bool
	Privileged      bool
}

type InstanceCreationOptions struct {
	RootPassword string
	UserName     string
	UserPassword string

	UserEmail    string
	UserFullName string
	SSHKeyURL    string

	AdditionalIPs     []string
	AdditionalDomains []string

	EnabledModules  []string
	EnabledServices []string
}

type Instance struct {
	Name   string
	Ports  []Port
	Status string
}

type Port struct {
	Service string
	Port    int32
}

type InstancesManager struct {
	docker *client.Client

	instancesMutex sync.Map
}

func NewInstancesManager(docker *client.Client) *InstancesManager {
	return &InstancesManager{
		docker: docker,
	}
}

// See https://stackoverflow.com/questions/64564781/golang-lock-per-value
func (m *InstancesManager) lockInstance(instanceName string) func() {
	value, _ := m.instancesMutex.LoadOrStore(instanceName, &sync.Mutex{})
	mtx := value.(*sync.Mutex)
	mtx.Lock()

	return func() { mtx.Unlock() }
}

func (m *InstancesManager) execInInstance(ctx context.Context, instanceName string, command []string) (string, string, int, error) {
	exec, err := m.docker.ContainerExecCreate(ctx, addContainerPrefix(instanceName), types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          command,
	})
	if err != nil {
		return "", "", 0, err
	}

	resp, err := m.docker.ContainerExecAttach(ctx, exec.ID, types.ExecStartCheck{})
	if err != nil {
		return "", "", 0, err
	}
	defer resp.Close()

	var outBuf, errBuf bytes.Buffer
	done := make(chan error)

	go func() {
		_, err := stdcopy.StdCopy(&outBuf, &errBuf, resp.Reader)

		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			return "", "", 0, err
		}
	case <-ctx.Done():
		return "", "", 0, ctx.Err()
	}

	stdout, err := io.ReadAll(&outBuf)
	if err != nil {
		return "", "", 0, err
	}

	stderr, err := io.ReadAll(&errBuf)
	if err != nil {
		return "", "", 0, err
	}

	meta, err := m.docker.ContainerExecInspect(ctx, exec.ID)
	if err != nil {
		return "", "", 0, err
	}

	return string(stdout), string(stderr), meta.ExitCode, nil
}

func (m *InstancesManager) execInInstanceStreaming(ctx context.Context, cancel func(error), instanceName string, command []string, stdoutChan, stderrChan chan []byte) (int, error) {
	rawStdoutChan, rawStderrChan := make(chan []byte), make(chan []byte)
	defer close(rawStdoutChan)
	defer close(rawStderrChan)

	exec, err := m.docker.ContainerExecCreate(ctx, addContainerPrefix(instanceName), types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          command,
	})
	if err != nil {
		cancel(err)

		return 0, err
	}

	resp, err := m.docker.ContainerExecAttach(ctx, exec.ID, types.ExecStartCheck{})
	if err != nil {
		cancel(err)

		return 0, err
	}

	done := make(chan struct{})
	go func() {
		for {
			header := make([]byte, 8)
			if n, err := resp.Reader.Read(header); err != nil {
				// EOF
				if n == 0 {
					done <- struct{}{}

					return
				}

				cancel(err)

				return
			}

			data := make([]byte, binary.BigEndian.Uint32(header[4:]))
			if _, err := resp.Reader.Read(data); err != nil {
				cancel(err)

				return
			}

			select {
			case <-ctx.Done():
				return
			default:
				switch header[0] {
				case 1:
					rawStderrChan <- data
				default:
					rawStdoutChan <- data
				}
			}
		}
	}()

	for {
		select {
		case <-done:
			meta, err := m.docker.ContainerExecInspect(ctx, exec.ID)
			if err != nil {
				return 0, err
			}

			return meta.ExitCode, err
		case <-ctx.Done():
			meta, err := m.docker.ContainerExecInspect(ctx, exec.ID)
			if err != nil {
				return 0, err
			}

			return meta.ExitCode, err
		case chunk := <-rawStderrChan:
			stdoutChan <- chunk
		case chunk := <-rawStdoutChan:
			stderrChan <- chunk
		}
	}
}

func (m *InstancesManager) GetInstances(ctx context.Context) ([]Instance, error) {
	containers, err := m.docker.ContainerList(ctx, types.ContainerListOptions{
		Filters: filters.NewArgs(filters.Arg("name", containerPrefix)),
		All:     true,
	})
	if err != nil {
		return []Instance{}, err
	}

	instances := []Instance{}
	for _, container := range containers {
		portsMap := make(map[int32]struct{})
		ports := []Port{}

		// Remove duplicates (i.e. IPv4 or IPv6)
		for _, port := range container.Ports {
			if _, exists := portsMap[int32(port.PublicPort)]; !exists {
				portsMap[int32(port.PublicPort)] = struct{}{}
				ports = append(ports, Port{
					Service: portToServiceMap[int32(port.PrivatePort)],
					Port:    int32(port.PublicPort),
				})
			}
		}

		instances = append(instances, Instance{
			Name:   removePrefix(container.Names[0]),
			Ports:  ports,
			Status: container.State,
		})
	}

	return instances, nil
}

func (m *InstancesManager) getFileContentsFromInstance(ctx context.Context, instanceName string, path string) (string, error) {
	// Copy the config file from the container
	// We use `path.Join` instead of `filepath.Join` here as the script at the path will always be executed in the container, which always uses `/` even if the host uses a different separator
	r, _, err := m.docker.CopyFromContainer(ctx, addContainerPrefix(instanceName), path)
	if err != nil {
		return "", err
	}

	// Read the config file into a buffer
	tr := tar.NewReader(r)
	var buf bytes.Buffer

	if _, err := tr.Next(); err != nil {
		return "", err
	}

	if _, err := io.Copy(&buf, tr); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (m *InstancesManager) writeFileToInstance(ctx context.Context, instanceName string, dest string, contents string, appendOnly bool) error {
	// If we are appending, get the old file contents
	if appendOnly {
		oldContents, err := m.getFileContentsFromInstance(ctx, instanceName, dest)
		if err != nil {
			return err
		}

		contents = oldContents + contents
	}

	// Create a tar archive containing the contents
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	if err := tw.WriteHeader(&tar.Header{
		Name: path.Base(dest),
		Mode: 664,
		Size: int64(len(contents)),
	}); err != nil {
		return err
	}
	if _, err := tw.Write([]byte(contents)); err != nil {
		return err
	}

	if err := tw.Close(); err != nil {
		return err
	}

	// Copy the file to the container
	return m.docker.CopyToContainer(ctx, addContainerPrefix(instanceName), path.Dir(dest), &buf, types.CopyToContainerOptions{})
}

func (m *InstancesManager) GetLogs(ctx context.Context, cancel func(error), instanceName string, stdoutChan, stderrChan chan []byte) {
	exitCode, err := m.execInInstanceStreaming(ctx, cancel, instanceName, journalctlTailCmd, stdoutChan, stderrChan)
	if err != nil {
		cancel(err)

		return
	}

	if exitCode != 0 {
		cancel(fmt.Errorf("could not get logs from instance: Non-zero exit code %v", exitCode))

		return
	}
}

func (m *InstancesManager) StartInstance(ctx context.Context, instanceName string, skipLock bool) error {
	if !skipLock {
		unlock := m.lockInstance(instanceName)
		defer unlock()
	}

	return m.docker.ContainerStart(ctx, addContainerPrefix(instanceName), types.ContainerStartOptions{})
}

func (m *InstancesManager) StopInstance(ctx context.Context, instanceName string) error {
	unlock := m.lockInstance(instanceName)
	defer unlock()

	return m.docker.ContainerStop(ctx, addContainerPrefix(instanceName), nil)
}

func (m *InstancesManager) RestartInstance(ctx context.Context, instanceName string) error {
	unlock := m.lockInstance(instanceName)
	defer unlock()

	return m.docker.ContainerRestart(ctx, addContainerPrefix(instanceName), nil)
}

func (m *InstancesManager) RemoveInstance(ctx context.Context, instanceName string, options InstanceRemovalOptions, skipLock bool) error {
	if !skipLock {
		unlock := m.lockInstance(instanceName)
		defer unlock()
	}

	// Remove customization; we have to call this before removing the container as it runs inside of it
	if options.Customizations {
		stdout, stderr, exitCode, err := m.execInInstance(ctx, instanceName, enumerateConfigurationScriptsCmd)
		if err != nil {
			return err
		}

		if exitCode != 0 {
			return fmt.Errorf("could not enumerate configuration scripts in instance: Non-zero exit code %v: stdout=%v, stderr=%v", exitCode, stdout, stderr)
		}

	o:
		for _, script := range strings.Split(string(stdout), "\n") {
			// Skip non-scripts
			if !strings.HasSuffix(script, ".sh") {
				continue
			}

			// Skip blocklisted scripts (such as the interactive parameters script)
			for _, blocklistedScript := range blocklistedScripts {
				if strings.Contains(script, blocklistedScript) {
					// We are in a nested loop, so continue at the outer one
					continue o
				}
			}

			// We use `path.Join` instead of `filepath.Join` here as the script at the path will always be executed in the container, which always uses `/` even if the host uses a different separator
			scriptPath := path.Join(configurationScriptsDir, script)

			stdout, stderr, exitCode, err := m.execInInstance(ctx, instanceName, getCommandToExecuteRefreshScript(scriptPath))
			if err != nil {
				return err
			}

			if exitCode != 0 {
				return fmt.Errorf("could not run configuration scripts in instance: Non-zero exit code %v: stdout=%v, stderr=%v", exitCode, stdout, stderr)
			}
		}
	}

	// Remove container
	if err := m.docker.ContainerRemove(ctx, addContainerPrefix(instanceName), types.ContainerRemoveOptions{
		Force: true,
	}); err != nil {
		return err
	}

	// Remove DEB cache
	if options.DEBCache {
		if err := m.docker.VolumeRemove(ctx, getDEBCacheVolumeName(instanceName), false); err != nil {
			return err
		}
	}

	// Remove preferences
	if options.Preferences {
		if err := m.docker.VolumeRemove(ctx, getPreferencesVolumeName(instanceName), false); err != nil {
			return err
		}
	}

	// Remove user data
	if options.UserData {
		for _, volume := range getUserDataVolumeNames(instanceName) {
			if err := m.docker.VolumeRemove(ctx, volume, false); err != nil {
				return err
			}
		}
	}

	// Remove transfer directory
	if options.Transfer {
		dir, err := getTransferDirectory(instanceName)
		if err != nil {
			return err
		}

		if err := os.RemoveAll(dir); err != nil {
			return err
		}
	}

	return nil
}

func (m *InstancesManager) GetCACert(ctx context.Context) (string, error) {
	volume, err := m.docker.VolumeInspect(ctx, getCAVolumeName())
	if err != nil {
		return "", err
	}

	cert, err := os.ReadFile(filepath.Join(volume.Mountpoint, caCertFile))
	if err != nil {
		return "", err
	}

	return string(cert), err
}

func (m *InstancesManager) ResetCA(ctx context.Context) error {
	return m.docker.VolumeRemove(ctx, getCAVolumeName(), false)
}

func (m *InstancesManager) GetShell(ctx context.Context, cancel func(error), instanceName string, stdinChan, stdoutChan, stderrChan chan []byte) {
	rawStdoutChan, rawStderrChan := make(chan []byte), make(chan []byte)
	defer close(rawStdoutChan)
	defer close(rawStderrChan)

	exec, err := m.docker.ContainerExecCreate(ctx, addContainerPrefix(instanceName), types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		AttachStdin:  true,
		Tty:          true,
		Cmd:          bashCmd,
	})
	if err != nil {
		cancel(err)

		return
	}

	resp, err := m.docker.ContainerExecAttach(ctx, exec.ID, types.ExecStartCheck{})
	if err != nil {
		cancel(err)

		return
	}

	go func() {
		for {
			header := make([]byte, 8)
			if _, err := resp.Reader.Read(header); err != nil {
				cancel(err)

				return
			}

			data := make([]byte, binary.BigEndian.Uint32(header[4:]))
			if _, err := resp.Reader.Read(data); err != nil {
				cancel(err)

				return
			}

			select {
			case <-ctx.Done():
				return
			default:
				switch header[0] {
				case 1:
					rawStderrChan <- data
				default:
					rawStdoutChan <- data
				}
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case chunk := <-stdinChan:
			if _, err := resp.Conn.Write(chunk); err != nil {
				cancel(err)

				return
			}
		case chunk := <-rawStderrChan:
			stdoutChan <- chunk
		case chunk := <-rawStdoutChan:
			stderrChan <- chunk
		}
	}
}

func (m *InstancesManager) ApplyInstance(ctx context.Context, cancel func(error), instanceName string, stdoutChan, stderrChan chan []byte, statusChan chan string, flags InstanceCreationFlags, opts InstanceCreationOptions) {
	unlock := m.lockInstance(instanceName)
	defer unlock()

	// Pull the latest image if specified
	if flags.PullLatestImage {
		statusChan <- getStartingMessage(StatusPullingLatestImage)

		if _, err := m.docker.ImagePull(ctx, pojdeDockerImage, types.ImagePullOptions{}); err != nil {
			cancel(err)

			return
		}

		statusChan <- getDoneMessage(StatusPullingLatestImage)
	}

	// Check if the container exists
	exists := true
	if _, err := m.docker.ContainerInspect(ctx, addContainerPrefix(instanceName)); err != nil {
		exists = false
	}

	// Remove the container if specified
	if flags.Recreate && exists {
		statusChan <- getStartingMessage(StatusRemovingExistingContainer)

		if err := m.RemoveInstance(ctx, instanceName, InstanceRemovalOptions{
			Customizations: false,
			DEBCache:       false,
			Preferences:    false,
			UserData:       false,
			Transfer:       false,
		}, true); err != nil {
			cancel(err)

			return
		}

		exists = false

		statusChan <- getDoneMessage(StatusRemovingExistingContainer)
	}

	// Create the container if it doesn't alreay exist
	if !exists {
		statusChan <- getStartingMessage(StatusCreatingContainer)

		hostConfig := &container.HostConfig{
			Mounts: []mount.Mount{},
		}

		// Allow access to Docker daemon from container
		if !flags.Isolate {
			hostConfig.Mounts = append(hostConfig.Mounts, mount.Mount{
				Type:   mount.TypeBind,
				Source: "/var/run/docker.sock",
				Target: "/var/run/docker.sock",
			})
		}

		// Enable privileged mode
		if flags.Privileged {
			hostConfig.Privileged = true
		}

		containerConfig := &container.Config{
			Env:     []string{},
			Volumes: map[string]struct{}{},
		}

		// Enable systemd
		containerConfig.Env = append(containerConfig.Env, "container=oci")

		hostConfig.Mounts = append(hostConfig.Mounts, mount.Mount{
			Type:         mount.TypeTmpfs,
			Target:       "/tmp",
			TmpfsOptions: &mount.TmpfsOptions{Mode: fs.FileMode(execPerm)},
		})

		hostConfig.Mounts = append(hostConfig.Mounts, mount.Mount{
			Type:         mount.TypeTmpfs,
			Target:       "/run",
			TmpfsOptions: &mount.TmpfsOptions{Mode: fs.FileMode(execPerm)},
		})

		hostConfig.Mounts = append(hostConfig.Mounts, mount.Mount{
			Type:         mount.TypeTmpfs,
			Target:       "/run/lock",
			TmpfsOptions: &mount.TmpfsOptions{Mode: fs.FileMode(execPerm)},
		})

		hostConfig.Mounts = append(hostConfig.Mounts, mount.Mount{
			Type:     mount.TypeBind,
			Source:   "/sys/fs/cgroup",
			Target:   "/sys/fs/cgroup",
			ReadOnly: true,
		})

		// Add preferences
		hostConfig.Mounts = append(hostConfig.Mounts, mount.Mount{
			Type:   mount.TypeVolume,
			Source: getPreferencesVolumeName(instanceName),
			Target: "/opt/pojde/preferences",
		})

		// Add CA volume
		hostConfig.Mounts = append(hostConfig.Mounts, mount.Mount{
			Type:   mount.TypeVolume,
			Source: getCAVolumeName(),
			Target: "/opt/pojde/ca",
		})

		// Add user data
		userDataVolumes := getUserDataVolumeNames(instanceName)

		hostConfig.Mounts = append(hostConfig.Mounts, mount.Mount{
			Type:   mount.TypeVolume,
			Source: userDataVolumes[0],
			Target: "/root",
		})

		hostConfig.Mounts = append(hostConfig.Mounts, mount.Mount{
			Type:   mount.TypeVolume,
			Source: userDataVolumes[1],
			Target: "/home",
		})

		// Add DEB cache
		hostConfig.Mounts = append(hostConfig.Mounts, mount.Mount{
			Type:   mount.TypeVolume,
			Source: getDEBCacheVolumeName(instanceName),
			Target: "/var/cache/apt/archives",
		})

		// Add transfer directory
		transfer, err := getTransferDirectory(instanceName)
		if err != nil {
			cancel(err)

			return
		}

		hostConfig.Mounts = append(hostConfig.Mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: transfer,
			Target: "/transfer",
		})

		containerConfig.Cmd = strslice.StrSlice{startCmd}
		containerConfig.Image = pojdeDockerImage

		// Expose ports
		exposedPorts := nat.PortSet{}
		portBindings := nat.PortMap{}

		for offset := 0; offset < portCount; offset++ {
			rawInternalPort := firstInternalPort + offset
			rawExternalPort := int(flags.StartPort) + offset

			containerPort := nat.Port(fmt.Sprintf("%v/tcp", rawInternalPort))

			exposedPorts[containerPort] = struct{}{}
			portBindings[containerPort] = []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: strconv.Itoa(rawExternalPort),
				},
				{
					HostIP:   "::",
					HostPort: strconv.Itoa(rawExternalPort),
				},
			}
		}

		containerConfig.ExposedPorts = exposedPorts
		hostConfig.PortBindings = portBindings

		// Always restart the container
		hostConfig.RestartPolicy = container.RestartPolicy{
			Name: "always",
		}

		// Create the container
		if _, err := m.docker.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, addContainerPrefix(instanceName)); err != nil {
			cancel(err)

			return
		}

		// Start the container
		if err := m.StartInstance(ctx, instanceName, true); err != nil {
			cancel(err)

			return
		}

		statusChan <- getDoneMessage(StatusCreatingContainer)

		statusChan <- getStartingMessage(StatusPreparingConfigScripts)

		// Prepare the config file
		configFile := config.NewConfig()

		configFile.RootPassword = opts.RootPassword
		configFile.UserName = opts.UserName
		configFile.UserPassword = opts.UserPassword

		configFile.UserEmail = opts.UserEmail
		configFile.UserFullName = opts.UserFullName
		configFile.SSHKeyURL = opts.SSHKeyURL

		configFile.AdditionalIPs = opts.AdditionalIPs
		configFile.AdditionalDomains = opts.AdditionalDomains

		configFile.EnabledModules = opts.EnabledModules
		configFile.EnabledServices = opts.EnabledServices

		configContents := configFile.Marshal()

		// Copy the config file to the container
		// We use `path.Join` instead of `filepath.Join` here it will always be used in the container, which always uses `/` even if the host uses a different separator
		if err := m.writeFileToInstance(ctx, instanceName, path.Join(preferencesDirInContainer, preferencesFileInContainer), configContents, false); err != nil {
			cancel(err)

			return
		}

		statusChan <- getDoneMessage(StatusPreparingConfigScripts)

		// Run the config scripts
		for _, script := range configScripts {
			statusChan <- getStartingMessage(GetConfigScriptStatus(script))

			// We use `path.Join` instead of `filepath.Join` here as the script at the path will always be executed in the container, which always uses `/` even if the host uses a different separator
			scriptPath := path.Join(configurationScriptsDir, script+".sh")

			exitCode, err := m.execInInstanceStreaming(ctx, cancel, instanceName, getCommandToExecuteUpgradeScript(scriptPath), stdoutChan, stderrChan)
			if err != nil {
				cancel(err)

				return
			}

			if exitCode != 0 {
				cancel(fmt.Errorf("could not run configuration scripts in instance: Non-zero exit code %v", exitCode))

				return
			}

			statusChan <- getDoneMessage(GetConfigScriptStatus(script))
		}

		cancel(nil)
	}
}

func (m *InstancesManager) GetInstanceConfig(ctx context.Context, instanceName string) (InstanceCreationOptions, error) {
	// Get the config file content
	cfgFileContent, err := m.getFileContentsFromInstance(ctx, instanceName, path.Join(preferencesDirInContainer, preferencesFileInContainer))
	if err != nil {
		return InstanceCreationOptions{}, err
	}

	// Parse the config file
	cfg := config.NewConfig()
	if err := cfg.Unmarshal(cfgFileContent); err != nil {
		return InstanceCreationOptions{}, err
	}

	return InstanceCreationOptions{
		RootPassword: cfg.RootPassword,
		UserName:     cfg.UserEmail,
		UserPassword: cfg.UserPassword,

		UserEmail:    cfg.UserEmail,
		UserFullName: cfg.UserFullName,
		SSHKeyURL:    cfg.SSHKeyURL,

		AdditionalIPs:     cfg.AdditionalIPs,
		AdditionalDomains: cfg.AdditionalDomains,

		EnabledModules:  cfg.EnabledModules,
		EnabledServices: cfg.EnabledServices,
	}, nil
}

func (m *InstancesManager) GetSSHKeys(ctx context.Context, instanceName string) ([]string, error) {
	// Get the authorized keys file
	authorizedKeysFileContent, err := m.getFileContentsFromInstance(ctx, instanceName, sshKeysPath)
	if err != nil {
		return []string{}, err
	}

	// Parse the SSH keys
	sshKeys := []string{}
	for _, sshKey := range strings.Split(authorizedKeysFileContent, "\n") {
		if !(strings.TrimSpace(sshKey) == "") {
			sshKeys = append(sshKeys, sshKey)
		}
	}

	// Parse the authorized keys
	return sshKeys, nil
}

func (m *InstancesManager) AddSSHKey(ctx context.Context, instanceName string, sshKey string) error {
	unlock := m.lockInstance(instanceName)
	defer unlock()

	// Add the key to the authorized_keys file
	return m.writeFileToInstance(ctx, instanceName, sshKeysPath, "\n"+sshKey, true)
}

func (m *InstancesManager) RemoveSSHKey(ctx context.Context, instanceName string, hash string) (string, error) {
	unlock := m.lockInstance(instanceName)
	defer unlock()

	// Get the current SSH keys
	sshKeys, err := m.GetSSHKeys(ctx, instanceName)
	if err != nil {
		return "", err
	}

	// Find all keys which don't match the hash
	targetKey := ""
	newSSHKeys := ""
	for _, sshKey := range sshKeys {
		if util.GetSHA512Hash(sshKey) != hash {
			newSSHKeys += sshKey + "\n"
		} else {
			targetKey = sshKey
		}
	}

	// If no key could be found for the target hash, return an error
	if targetKey == "" {
		return "", fmt.Errorf("could not find SSH key for hash %v", hash)
	}

	// Write new keys
	if err := m.writeFileToInstance(ctx, instanceName, sshKeysPath, newSSHKeys, false); err != nil {
		return "", err
	}

	return targetKey, nil
}
