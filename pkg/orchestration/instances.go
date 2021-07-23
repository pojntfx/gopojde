package orchestration

import (
	"context"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

const (
	prefix = "/pojde-"
)

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

func (m *InstancesManager) GetInstances() ([]Instance, error) {
	containers, err := m.docker.ContainerList(context.Background(), types.ContainerListOptions{
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
			Name:   strings.TrimPrefix(container.Names[0], prefix),
			Ports:  ports,
			Status: container.State,
		})
	}

	return instances, nil
}
