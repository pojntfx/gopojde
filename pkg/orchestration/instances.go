package orchestration

import (
	"context"
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

const (
	prefix = "/pojde-"
)

func removePrefix(name string) string {
	return strings.TrimPrefix(name, prefix)
}

func addPrefix(name string) string {
	return prefix + name
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

func (m *InstancesManager) GetLogs(ctx context.Context, instanceName string) (chan string, chan error) {
	outChan := make(chan string)
	errChan := make(chan error)

	logs, err := m.docker.ContainerLogs(ctx, addPrefix(instanceName), types.ContainerLogsOptions{
		ShowStderr: true,
		ShowStdout: true,
		Timestamps: false,
		Follow:     true,
	})
	if err != nil {
		// We have to launch this in a new Goroutine to prevent a deadlock as the write operation to the chan would be blocking
		go func() {
			errChan <- fmt.Errorf("could not get instance logs: %v", err)

			close(outChan)
			close(errChan)
		}()

		return outChan, errChan
	}

	go func() {
		header := make([]byte, 8)
		for {
			if _, err := logs.Read(header); err != nil {
				errChan <- fmt.Errorf("could not read from instance log stream: %v", err)

				close(outChan)
				close(errChan)

				return
			}

			data := make([]byte, binary.BigEndian.Uint32(header[4:]))
			if _, err := logs.Read(data); err != nil {
				errChan <- fmt.Errorf("could not read from instance log stream: %v", err)

				close(outChan)
				close(errChan)

				return
			}

			outChan <- string(data)
		}
	}()

	return outChan, nil
}
