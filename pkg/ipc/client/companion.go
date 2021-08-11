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

func (c *CompanionIPCClient) CreateSSHConnection(
	instanceID string,
	privateKey string,
	// TODO: Prevent this callback from being 0
	passwordGetterFunc func() string,
	hostKeyValidatorFunc func(hostname, fingerprint string) error,
) (string, error) {
	passwordGetter := app.FuncOf(func(this app.Value, args []app.Value) interface{} {
		return passwordGetterFunc()
	})
	defer passwordGetter.Release()

	hostKeyValidator := app.FuncOf(func(this app.Value, args []app.Value) interface{} {
		return hostKeyValidatorFunc(args[0].String(), args[1].String())
	})
	defer hostKeyValidator.Release()

	res, err := interop.Await(app.Window().Call(
		shared.CreateSSHConnection,
		instanceID,
		privateKey,
		passwordGetter,
		hostKeyValidator,
	))
	if err != nil {
		return "", err
	}

	return res.String(), err
}
