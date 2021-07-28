package orchestration

import (
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

	"github.com/alessio/shellescape"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
)

const (
	containerPrefix         = "/pojde-"
	volumePrefix            = "pojde-"
	configurationScriptsDir = "/opt/pojde/configuration"
	caCertFile              = "ca.pem"
	pojdeDockerImage        = "pojntfx/pojde:latest"
	firstPort               = 8000
	portCount               = 5
	startCmd                = "/lib/systemd/systemd"
	execPerm                = 0775
)

var (
	journalctlTailCmd                = []string{"journalctl", "-f"}
	enumerateConfigurationScriptsCmd = []string{"ls", configurationScriptsDir}
	bashCmd                          = []string{"bash"}
	blocklistedScripts               = []string{"parameters.sh"} // TODO: Remove this as soon as scripts are no longer interactive
)

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

type InstanceRemovalOptions struct {
	Customizations bool
	DEBCache       bool
	Preferences    bool
	UserData       bool
	Transfer       bool
}

type InstanceCreationFlags struct {
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
	SSHKey       string

	AdditionalIPs     []string
	AdditionalDomains []string

	EnabledModules  []string
	EnabledServices []string
}

type Instance struct {
	Name   string
	Ports  []int32
	Status string
}

type InstancesManager struct {
	docker *client.Client
}

func NewInstancesManager(docker *client.Client) *InstancesManager {
	return &InstancesManager{
		docker: docker,
	}
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
		ports := []int32{}

		// Remove duplicates (i.e. IPv4 or IPv6)
		for _, port := range container.Ports {
			if _, exists := portsMap[int32(port.PublicPort)]; !exists {
				portsMap[int32(port.PublicPort)] = struct{}{}
				ports = append(ports, int32(port.PublicPort))
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

func (m *InstancesManager) GetLogs(ctx context.Context, instanceName string) (chan string, chan error, *types.HijackedResponse) {
	outChan := make(chan string)
	errChan := make(chan error)

	exec, err := m.docker.ContainerExecCreate(ctx, addContainerPrefix(instanceName), types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          journalctlTailCmd,
	})
	if err != nil {
		// We have to launch this in a new Goroutine to prevent a deadlock as the write operation to the chan would be blocking
		go func() {
			errChan <- fmt.Errorf("could not request instance logs: %v", err)

			close(outChan)
			close(errChan)
		}()

		return outChan, errChan, nil
	}

	logs, err := m.docker.ContainerExecAttach(ctx, exec.ID, types.ExecStartCheck{})
	if err != nil {
		// We have to launch this in a new Goroutine to prevent a deadlock as the write operation to the chan would be blocking
		go func() {
			errChan <- fmt.Errorf("could not attach to instance logs: %v", err)

			close(outChan)
			close(errChan)
		}()

		return outChan, errChan, &logs
	}

	go func() {
		for {
			header := make([]byte, 8)
			if _, err := logs.Reader.Read(header); err != nil {
				errChan <- fmt.Errorf("could not read from instance log stream: %v", err)

				close(outChan)
				close(errChan)

				return
			}

			data := make([]byte, binary.BigEndian.Uint32(header[4:]))
			if _, err := logs.Reader.Read(data); err != nil {
				errChan <- fmt.Errorf("could not read from instance log stream: %v", err)

				close(outChan)
				close(errChan)

				return
			}

			outChan <- string(data)
		}
	}()

	return outChan, nil, &logs
}

func (m *InstancesManager) StartInstance(ctx context.Context, instanceName string) error {
	return m.docker.ContainerStart(ctx, addContainerPrefix(instanceName), types.ContainerStartOptions{})
}

func (m *InstancesManager) StopInstance(ctx context.Context, instanceName string) error {
	return m.docker.ContainerStop(ctx, addContainerPrefix(instanceName), nil)
}

func (m *InstancesManager) RestartInstance(ctx context.Context, instanceName string) error {
	return m.docker.ContainerRestart(ctx, addContainerPrefix(instanceName), nil)
}

func (m *InstancesManager) RemoveInstance(ctx context.Context, instanceName string, options InstanceRemovalOptions) error {
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

func (m *InstancesManager) ApplyInstance(ctx context.Context, name string, flags InstanceCreationFlags, opts InstanceCreationOptions) error {
	// Pull the latest image if specified
	if flags.PullLatestImage {
		if _, err := m.docker.ImagePull(ctx, pojdeDockerImage, types.ImagePullOptions{}); err != nil {
			return err
		}
	}

	// Check if the container exists
	exists := true
	if _, err := m.docker.ContainerInspect(ctx, addContainerPrefix(name)); err != nil {
		exists = false
	}

	// Remove the container if specified
	if flags.Recreate && exists {
		if err := m.RemoveInstance(ctx, name, InstanceRemovalOptions{
			Customizations: false,
			DEBCache:       false,
			Preferences:    false,
			UserData:       false,
			Transfer:       false,
		}); err != nil {
			return err
		}

		exists = false
	}

	// Create the container if it doesn't alreay exist
	if !exists {
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
			Source: getPreferencesVolumeName(name),
			Target: "/opt/pojde/preferences",
		})

		// Add CA volume
		hostConfig.Mounts = append(hostConfig.Mounts, mount.Mount{
			Type:   mount.TypeVolume,
			Source: getCAVolumeName(),
			Target: "/opt/pojde/ca",
		})

		// Add user data
		userDataVolumes := getUserDataVolumeNames(name)

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
			Source: getDEBCacheVolumeName(name),
			Target: "/var/cache/apt/archives",
		})

		// Add transfer directory
		transfer, err := getTransferDirectory(name)
		if err != nil {
			return err
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
			rawPort := firstPort + offset
			containerPort := nat.Port(fmt.Sprintf("%v/tcp", rawPort))

			exposedPorts[containerPort] = struct{}{}
			portBindings[containerPort] = []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: strconv.Itoa(rawPort),
				},
				{
					HostIP:   "::",
					HostPort: strconv.Itoa(rawPort),
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
		resp, err := m.docker.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, addContainerPrefix(name))
		if err != nil {
			return err
		}

		if err := m.docker.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
			return err
		}
	}

	return nil
}
