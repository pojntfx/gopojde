package components

import "github.com/maxence-charriere/go-app/v9/pkg/app"

type CompanionHome struct {
	app.Compo
}

func (c *CompanionHome) Render() app.UI {
	return app.H1().Text("Hello, gopojde Companion!")
}
