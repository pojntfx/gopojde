package services

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	api "github.com/pojntfx/gopojde/pkg/api/proto/v1"
	"github.com/pojntfx/gopojde/pkg/orchestration"
	"google.golang.org/protobuf/types/known/emptypb"
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

func (s *InstancesService) GetInstances(ctx context.Context, _ *empty.Empty) (*api.InstancesMessage, error) {
	instances, err := s.instancesManager.GetInstances(ctx)
	if err != nil {
		return &api.InstancesMessage{}, err
	}

	out := []*api.InstanceMessage{}
	for _, instance := range instances {
		out = append(out, &api.InstanceMessage{
			Name:   instance.Name,
			Ports:  instance.Ports,
			Status: instance.Status,
		})
	}

	return &api.InstancesMessage{
		Instances: out,
	}, nil
}

func (s *InstancesService) GetLogs(req *api.InstanceReferenceMessage, stream api.InstancesService_GetLogsServer) error {
	logChan, errChan, handle := s.instancesManager.GetLogs(stream.Context(), req.GetName())
	defer handle.Close()

	for {
		select {
		case chunk := <-logChan:
			if err := stream.Send(&api.LogMessage{
				Chunk: chunk,
			}); err != nil {
				return err
			}
		case err := <-errChan:
			return err
		}
	}
}

func (s *InstancesService) StartInstance(ctx context.Context, req *api.InstanceReferenceMessage) (*empty.Empty, error) {
	return &emptypb.Empty{}, s.instancesManager.StartInstance(ctx, req.GetName())
}

func (s *InstancesService) StopInstance(ctx context.Context, req *api.InstanceReferenceMessage) (*empty.Empty, error) {
	return &emptypb.Empty{}, s.instancesManager.StopInstance(ctx, req.GetName())
}

func (s *InstancesService) RestartInstance(ctx context.Context, req *api.InstanceReferenceMessage) (*empty.Empty, error) {
	return &emptypb.Empty{}, s.instancesManager.RestartInstance(ctx, req.GetName())
}
