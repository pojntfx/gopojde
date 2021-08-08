package client

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/pojntfx/gopojde/pkg/interop"
	"github.com/pojntfx/gopojde/pkg/ipc/shared"
)

type CompanionIPCClient struct {
}

func NewCompanionIPCClient() *CompanionIPCClient {
	return &CompanionIPCClient{}
}

func (c *CompanionIPCClient) Open(ctx app.Context, address string) error {
	if _, err := interop.Await(app.Window().Call(shared.OpenKey, address)); err != nil {
		return err
	}

	return nil
}

func (c *CompanionIPCClient) GetInstances() ([]shared.Instance, error) {
	res, err := interop.Await(app.Window().Call(shared.GetInstancesKey))
	if err != nil {
		return []shared.Instance{}, err
	}

	rv := []shared.Instance{}
	if err := interop.Unmarshal(res, &rv); err != nil {
		return []shared.Instance{}, err
	}

	return rv, nil
}
