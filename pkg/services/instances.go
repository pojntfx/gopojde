package services

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	api "github.com/pojntfx/gopojde/pkg/api/proto/v1"
	"github.com/pojntfx/gopojde/pkg/orchestration"
)

//go:generate sh -c "mkdir -p ../api/proto/v1 && protoc --go_out=paths=source_relative,plugins=grpc:../api/proto/v1 -I=../../api/proto/v1 ../../api/proto/v1/*.proto"

type InstancesService struct {
	api.UnimplementedInstancesServiceServer

	instancesManager *orchestration.InstancesManager
}

func NewInstancesService(instancesManager *orchestration.InstancesManager) *InstancesService {
	return &InstancesService{
		instancesManager: instancesManager,
	}
}

func (s *InstancesService) GetInstances(context.Context, *empty.Empty) (*api.Instances, error) {
	instances, err := s.instancesManager.GetInstances()
	if err != nil {
		return &api.Instances{}, err
	}

	out := []*api.Instance{}
	for _, instance := range instances {
		out = append(out, &api.Instance{
			Name:   instance.Name,
			Ports:  instance.Ports,
			Status: instance.Status,
		})
	}

	return &api.Instances{
		Instances: out,
	}, nil
}
