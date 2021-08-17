package components

import (
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

func (c *CompanionHome) refresh(ctx app.Context) error {
	if !c.connected {
		if err := c.ipc.Open(ctx, "ws://localhost:15324"); err != nil {
			return err
		}

		c.connected = true
	}

	instances, err := c.ipc.GetInstances(app.Window().Call("prompt", "SSH private key").String())
	if err != nil {
		return err
	}

	c.instances = instances

	return nil
}

func (c *CompanionHome) Render() app.UI {
	return app.Div().Class("pf-c-page").Body(
		app.Header().Class("pf-c-page__header").Body(
			// TODO: Use gopojde logo
			app.Div().Class("pf-c-page__header-brand").Body(
				app.A().Class("pf-c-page__header-brand-link").Body(
					app.Img().Class("pf-c-brand").Src("/assets/images/PF-Masthead-Logo.svg").Alt("PatternFly logo"),
				),
			),
			app.Div().Class("pf-c-page__header-tools").Body(
				app.Div().Class("pf-c-page__header-tools-group").Body(
					app.Div().Class("pf-c-page__header-tools-item").Body(
						app.Button().Class("pf-c-button pf-m-plain").Type("button").Aria("label", "Refresh").OnClick(func(ctx app.Context, e app.Event) {
							if err := c.refresh(ctx); err != nil {
								log.Fatal(err)
							}
						}).Body(
							app.I().Class("fas fa-sync").Aria("hidden", true),
						),
					),
				),
			),
		),
		app.Main().Class("pf-c-page__main").TabIndex(-1).Body(
			app.Section().Class("pf-c-page__main-section pf-m-no-padding pf-m-padding-on-xl").Body(
				app.Div().Class("pf-c-card").Body(
					app.Ul().Class("pf-c-data-list pf-x-u-border-top-0").Aria("role", "list").Aria("label", "List of instances").Body(
						app.Range(c.instances).Slice(func(i int) app.UI {
							return app.Li().Class("pf-c-data-list__item").Body(
								app.Div().Class("pf-c-data-list__item-row").Body(
									app.Div().Class("pf-c-data-list__item-content").Body(
										app.Div().Class("pf-c-data-list__cell pf-m-align-left").Body(
											app.Div().Class("pf-l-flex pf-m-column pf-m-space-items-md").Body(
												app.Div().Class("pf-l-flex pf-m-column pf-m-space-items-none").Body(
													app.Div().Class("pf-l-flex__item").Body(
														app.P().Text(c.instances[i].Name),
													),
												),
												app.If(len(c.instances[i].Tunnels) > 0,
													app.Div().Class("pf-l-flex__item").Body(
														app.Div().Class("pf-c-chip-group").Body(
															app.Div().Class("pf-c-chip-group__main").Body(
																app.Ul().Class("pf-c-chip-group__list").Aria("role", "list").Aria("label", "Open ports").Body(
																	app.Li().Class("pf-c-chip-group__list-item").Body(
																		app.Range(c.instances[i].Tunnels).Slice(func(j int) app.UI {
																			return app.Div().Class("pf-c-chip").Body(
																				app.Span().Class("pf-c-chip__text").Text(c.instances[i].Tunnels[j].LocalAddress+"→"+c.instances[i].Tunnels[j].RemoteAddress),
																				app.Button().Class("pf-c-button pf-m-plain").Type("button").Aria("label", "Remove port").Body(
																					app.I().Class("fas fa-times").Aria("hidden", true),
																				),
																			)
																		}),
																	),
																),
															),
														),
													),
												),
											),
										),
										app.Div().Class("pf-c-data-list__cell pf-m-align-right pf-m-no-fill pf-u-mt-md-on-md").Body(
											app.Button().Class("pf-c-button pf-m-secondary").Type("button").Aria("label", "Add a port").OnClick(func(ctx app.Context, e app.Event) {
												if _, err := c.ipc.ForwardFromLocalToRemote(ctx, c.instances[i].ID, app.Window().Call("prompt", "Local address").String(), app.Window().Call("prompt", "Remote address").String()); err != nil {
													log.Fatal(err)
												}

												if err := c.refresh(ctx); err != nil {
													log.Fatal(err)
												}
											}).Body(
												app.I().Class("fas fa-plus").Aria("hidden", true),
											),
										),
									),
								),
							)
						}),
					),
				),
			),
		),
	)
}

func (c *CompanionHome) OnMount(ctx app.Context) {
	c.instances = []shared.Instance{}

	c.ipc = client.NewCompanionIPCClient()
	c.ipc.Bind(ctx)
}
