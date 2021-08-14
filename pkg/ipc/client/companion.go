package client

import (
	"errors"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/pojntfx/gopojde/pkg/interop"
	"github.com/pojntfx/gopojde/pkg/ipc/shared"
)

type CompanionIPCClient struct {
}

func NewCompanionIPCClient() *CompanionIPCClient {
	return &CompanionIPCClient{}
}

func (c *CompanionIPCClient) Bind(ctx app.Context) error {
	app.Window().Set(shared.PasswordGetterKey, app.FuncOf(func(this app.Value, args []app.Value) interface{} {
		return app.Window().Call("prompt", "SSH private key's password").String()
	}))

	app.Window().Set(shared.HostKeyValidatorKey, app.FuncOf(func(this app.Value, args []app.Value) interface{} {
		confirmed := app.Window().Call("confirm", `Does the fingerprint "`+args[1].String()+`" match for the hostname "`+args[0].String()+`"?`).Bool()
		if !confirmed {
			return errors.New("fingerprint did not match for host")
		}

		return nil
	}))

	return nil
}

func (c *CompanionIPCClient) Open(ctx app.Context, address string) error {
	if _, err := interop.Await(app.Window().Call(shared.OpenKey, address)); err != nil {
		return err
	}

	return nil
}

func (c *CompanionIPCClient) GetInstances(privateKey string) ([]shared.Instance, error) {
	res, err := interop.Await(app.Window().Call(shared.GetInstancesKey, privateKey))
	if err != nil {
		return []shared.Instance{}, err
	}

	rv := []shared.Instance{}
	if err := interop.Unmarshal(res, &rv); err != nil {
		return []shared.Instance{}, err
	}

	return rv, nil
}
