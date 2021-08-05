package components

import "github.com/maxence-charriere/go-app/v9/pkg/app"

type ManagerHome struct {
	app.Compo
}

func (c *ManagerHome) Render() app.UI {
	return app.H1().Text("Hello, gopojde Manager!")
}
