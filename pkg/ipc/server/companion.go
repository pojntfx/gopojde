package server

import (
	"context"
	"errors"
	"time"

	"github.com/pojntfx/go-app-grpc-chat-frontend-web/pkg/websocketproxy"
	api "github.com/pojntfx/gopojde/pkg/api/proto/v1"
	"github.com/pojntfx/gopojde/pkg/ipc/shared"
	"github.com/zserge/lorca"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CompanionIPCServer struct {
	daemon api.InstancesServiceClient
}

func NewCompanionIPC() *CompanionIPCServer {
	return &CompanionIPCServer{}
}

func (c *CompanionIPCServer) Bind(ui lorca.UI) error {
	if err := ui.Bind(shared.OpenKey, c.Open); err != nil {
		return err
	}

	if err := ui.Bind(shared.GetInstancesKey, c.GetInstances); err != nil {
		return err
	}

	return nil
}

func (c *CompanionIPCServer) Open(address string) error {
	conn, err := grpc.Dial(address, grpc.WithContextDialer(websocketproxy.NewWebSocketProxyClient(time.Minute).Dialer), grpc.WithInsecure())
	if err != nil {
		return err
	}

	c.daemon = api.NewInstancesServiceClient(conn)

	return nil
}

func (c *CompanionIPCServer) GetInstances() ([]shared.Instance, error) {
	if c.daemon == nil {
		return []shared.Instance{}, errors.New("could not get instances: not connected to daemon")
	}

	// Get all instances
	instances, err := c.daemon.GetInstances(context.Background(), &emptypb.Empty{})
	if err != nil {
		return []shared.Instance{}, err
	}

	// Reduce instances to relevant options
	res := []shared.Instance{}
	for _, instance := range instances.GetInstances() {
		res = append(res, shared.Instance{
			ID:   instance.GetInstanceID().GetName(),
			Name: instance.GetInstanceID().GetName(),
		})
	}

	return res, nil
}
