package components

import (
	"errors"
	"log"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/pojntfx/gopojde/pkg/ipc/client"
	"github.com/pojntfx/gopojde/pkg/ipc/shared"
)

type CompanionHome struct {
	app.Compo

	connected bool
	instances []shared.Instance

	ipc *client.CompanionIPCClient
}

func (c *CompanionHome) Render() app.UI {
	return app.Div().Body(
		app.H1().Class("pf-c-title").Text("gopojde Companion"),
		app.If(c.connected,
			app.Button().Class("pf-c-button pf-m-primary").Type("button").Text("Get instances").OnClick(func(ctx app.Context, e app.Event) {
				instances, err := c.ipc.GetInstances()
				if err != nil {
					log.Fatal(err)
				}

				c.instances = instances
			}),
			app.Ul().Class("pf-c-list").Body(
				app.Range(c.instances).Slice(func(i int) app.UI {
					return app.Li().Body(
						app.Span().Text(c.instances[i].Name),
						app.Button().Class("pf-c-button").Text("Connect via SSH").OnClick(func(ctx app.Context, e app.Event) {
							key, err := c.ipc.CreateSSHConnection(
								c.instances[i].ID,
								app.Window().Call("prompt", "SSH private key").String(),
								func() string {
									return app.Window().Call("prompt", "SSH private key's password").String()
								},
								func(hostname, fingerprint string) error {
									confirmed := app.Window().Call("confirm", `Does the fingerprint "`+fingerprint+`" match for the hostname "`+hostname+`"?`).Bool()
									if !confirmed {
										return errors.New("fingerprint did not match for hostname")
									}

									return nil
								},
							)
							if err != nil {
								log.Fatal(err)
							}

							log.Println("Created SSH connection with key", key)
						}),
					)
				}),
			),
		).Else(
			app.Button().Class("pf-c-button pf-m-primary").Type("button").Text("Connect to backend").OnClick(func(ctx app.Context, e app.Event) {
				if err := c.ipc.Open(ctx, "ws://localhost:15324"); err != nil {
					log.Fatal(err)
				}

				c.connected = true
			}),
		),
	)
}

func (c *CompanionHome) OnMount(app.Context) {
	c.instances = []shared.Instance{}

	c.ipc = client.NewCompanionIPCClient()
}
