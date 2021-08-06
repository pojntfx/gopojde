package components

import (
	"log"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type CompanionHome struct {
	app.Compo
}

func (c *CompanionHome) Render() app.UI {
	return app.Div().Body(
		app.H1().Text("gopojde Companion"),
		app.Button().Text("Connect to backend").OnClick(func(ctx app.Context, e app.Event) {
			rv := app.Window().Call("gopojdeCompanionConnectToDaemon", "ws://localhost:15324")

			log.Println(rv.JSValue().Type())
		}),
	)
}
