package orchestration

import (
	"context"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

const (
	prefix = "/pojde-"
)

var (
	journalctlTail = []string{"journalctl", "-f"}
)

func removePrefix(name string) string {
	return strings.TrimPrefix(name, prefix)
}

func addPrefix(name string) string {
	return prefix + name
}

func getDEBCacheVolumeName(name string) string {
	return addPrefix(name + "-apt-cache")
}

func getPreferencesVolumeName(name string) string {
	return addPrefix(name + "-preferences")
}

func getUserDataVolumeNames(name string) []string {
	return []string{addPrefix(name + "-home-root"), addPrefix(name + "-home-user")}
}

func getTransferDirectory(name string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, "Documents", "pojde", name), nil
}

type InstanceRemovalOptions struct {
	Customizations bool
	DEBCache       bool
	Preferences    bool
	UserData       bool
	Transfer       bool
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

func (m *InstancesManager) GetInstances(ctx context.Context) ([]Instance, error) {
	containers, err := m.docker.ContainerList(ctx, types.ContainerListOptions{
		Filters: filters.NewArgs(filters.Arg("name", prefix)),
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

	exec, err := m.docker.ContainerExecCreate(ctx, addPrefix(instanceName), types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          journalctlTail,
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
	return m.docker.ContainerStart(ctx, addPrefix(instanceName), types.ContainerStartOptions{})
}

func (m *InstancesManager) StopInstance(ctx context.Context, instanceName string) error {
	return m.docker.ContainerStop(ctx, addPrefix(instanceName), nil)
}

func (m *InstancesManager) RestartInstance(ctx context.Context, instanceName string) error {
	return m.docker.ContainerRestart(ctx, addPrefix(instanceName), nil)
}

func (m *InstancesManager) RemoveInstance(ctx context.Context, instanceName string, options InstanceRemovalOptions) error {
	// Remove customization; we have to call this before removing the container as it runs inside of it
	if options.Customizations {
		// TODO: Loop over vendored scripts and call `refresh`
	}

	// Remove container
	if err := m.docker.ContainerRemove(ctx, addPrefix(instanceName), types.ContainerRemoveOptions{
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
