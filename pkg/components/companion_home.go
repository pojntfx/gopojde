package components

import (
	"encoding/json"
	"log"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
	api "github.com/pojntfx/gopojde/pkg/api/proto/v1"
	"github.com/pojntfx/gopojde/pkg/interop"
)

type InstanceAndOptions struct {
	*api.InstanceSummaryMessage
	*api.InstanceOptionsMessage
}

type CompanionHome struct {
	app.Compo

	connected bool
	instances []InstanceAndOptions
}

func (c *CompanionHome) Render() app.UI {
	return app.Div().Body(
		app.H1().Class("pf-c-title").Text("gopojde Companion"),
		app.If(c.connected,
			app.Button().Class("pf-c-button pf-m-primary").Type("button").Text("Get instances").OnClick(func(ctx app.Context, e app.Event) {
				rv, err := interop.Await(app.Window().Call("gopojdeCompanionGetInstances"))
				if err != nil {
					log.Fatal(err)
				}

				if err := json.Unmarshal([]byte(app.Window().Get("JSON").Call("stringify", rv).String()), &c.instances); err != nil {
					log.Fatal(err)
				}
			}),
			app.Ul().Class("pf-c-list").Body(
				app.Range(c.instances).Slice(func(i int) app.UI {
					return app.Li().Text(c.instances[i])
				}),
			),
		).Else(
			app.Button().Class("pf-c-button pf-m-primary").Type("button").Text("Connect to backend").OnClick(func(ctx app.Context, e app.Event) {
				if _, err := interop.Await(app.Window().Call("gopojdeCompanionConnectToDaemon", "ws://localhost:15324")); err != nil {
					log.Fatal(err)
				}

				c.connected = true
			}),
		),
	)
}

func (c *CompanionHome) OnMount(app.Context) {
	c.instances = []InstanceAndOptions{}
}
