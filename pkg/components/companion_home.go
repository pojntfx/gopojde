package components

import (
	"log"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/pojntfx/gopojde/pkg/interop"
)

type CompanionHome struct {
	app.Compo

	connected bool
}

func (c *CompanionHome) Render() app.UI {
	return app.Div().Body(
		app.H1().Text("gopojde Companion"),
		app.If(c.connected,
			app.Button().Text("Get instances").OnClick(func(ctx app.Context, e app.Event) {
				instances, err := interop.Await(app.Window().Call("gopojdeCompanionGetInstances"))
				if err != nil {
					app.Window().Call("alert", err.Error())

					return
				}

				log.Println("got instances:", instances)
			}),
		).Else(
			app.Button().Text("Connect to backend").OnClick(func(ctx app.Context, e app.Event) {
				if _, err := interop.Await(app.Window().Call("gopojdeCompanionConnectToDaemon", "ws://localhost:15324")); err != nil {
					app.Window().Call("alert", err.Error())

					return
				}

				c.connected = true
			}),
		),
	)
}
