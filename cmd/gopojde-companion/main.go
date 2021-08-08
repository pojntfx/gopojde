package main

import (
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/pojntfx/gopojde/pkg/components"
	"github.com/pojntfx/gopojde/pkg/ipc/server"
	"github.com/pojntfx/gopojde/pkg/web/companion"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zserge/lorca"
)

const (
	buildKey = "build"
	outKey   = "out"
	pathKey  = "path"
	serveKey = "serve"
	laddrKey = "laddr"
)

func main() {
	// Web code
	{
		// Define the routes
		app.Route("/", &components.CompanionHome{})

		// Start the app
		app.RunWhenOnBrowser()
	}

	// Builder code
	{
		if os.Getenv("BUILDER") == "true" {
			// Create command
			cmd := &cobra.Command{
				Use:   "gopojde-companion",
				Short: "Experimental Go implementation of https://github.com/pojntfx/pojde.",
				Long: `Experimental Go implementation of https://github.com/pojntfx/pojde.

For more information, please visit https://github.com/pojntfx/gopojde.`,
				RunE: func(cmd *cobra.Command, args []string) error {
					// Define the handler
					h := &app.Handler{
						Author:          "Felicitas Pojtinger",
						BackgroundColor: "#151515",
						Description:     "Experimental Go implementation of https://github.com/pojntfx/pojde.",
						Icon: app.Icon{
							Default: "/web/icon.png",
						},
						Keywords: []string{
							"vscode",
							"novnc",
							"cockpit",
							"ttyd",
							"juypter-lab",
							"code-server",
						},
						LoadingLabel: "Experimental Go implementation of https://github.com/pojntfx/pojde.",
						Name:         "gopojde",
						RawHeaders: []string{
							`<meta property="og:url" content="https://pojntfx.github.io/gopojde/">`,
							`<meta property="og:title" content="gopojde">`,
							`<meta property="og:description" content="Experimental Go implementation of https://github.com/pojntfx/pojde.">`,
							`<meta property="og:image" content="https://pojntfx.github.io/gopojde/web/icon.png">`,
						},
						Styles: []string{
							`https://unpkg.com/@patternfly/patternfly@4.115.2/patternfly.css`,
							`https://unpkg.com/@patternfly/patternfly@4.115.2/patternfly-addons.css`,
							`/web/index.css`,
						},
						ThemeColor: "#151515",
						Title:      "gopojde",
					}

					// Create static build if specified
					if viper.GetBool(buildKey) {
						// Deploy under a path
						if path := viper.GetString(pathKey); path != "" {
							h.Resources = app.GitHubPages(path)
						}

						if err := app.GenerateStaticWebsite(viper.GetString(outKey), h); err != nil {
							return err
						}
					}

					// Serve if specified
					if viper.GetBool(serveKey) {
						log.Printf("gopojde companion listening on %v", viper.GetString(laddrKey))

						if err := http.ListenAndServe(viper.GetString(laddrKey), h); err != nil {
							return err
						}
					}

					return nil
				},
			}

			// Bind flags
			cmd.PersistentFlags().Bool(buildKey, false, "Create static build")
			cmd.PersistentFlags().String(outKey, "out/gopojde-companion-web", "Out directory for static build")
			cmd.PersistentFlags().String(pathKey, "", "Base path for static build")
			cmd.PersistentFlags().Bool(serveKey, false, "Build and serve the companion")
			cmd.PersistentFlags().String(laddrKey, "localhost:15326", "Address to serve the companion on")

			// Bind env variables
			if err := viper.BindPFlags(cmd.PersistentFlags()); err != nil {
				log.Fatal(err)
			}
			viper.SetEnvPrefix("gopojde_companion")
			viper.AutomaticEnv()

			// Run command
			if err := cmd.Execute(); err != nil {
				log.Fatal(err)
			}

			return
		}
	}

	// Native code
	{
		// Name the instance
		if runtime.GOOS == "linux" {
			os.Args = append(os.Args, "--class=gopojde Companion")
		}

		// Spawn Chromium instance
		ui, err := lorca.New("", "", 480, 640, os.Args...)
		if err != nil {
			log.Fatal("could not spawn Chromium instance:", err)
		}
		defer ui.Close()

		// Bind IPC handlers
		companionIPC := server.NewCompanionIPC()
		companionIPC.Bind(ui)

		// Start integrated webserver
		lis, err := net.Listen("tcp", "localhost:0")
		if err != nil {
			log.Fatal("could not listen:", err)
		}
		defer lis.Close()

		// Serve the compiled companion
		root, err := fs.Sub(companion.FS, "assets")
		if err != nil {
			log.Fatal("could not get embedded root filesystem:", err)
		}

		go http.Serve(lis, http.FileServer(http.FS(root)))

		// Open the page on the integrated webserver
		if err := ui.Load(fmt.Sprintf("http://%v", lis.Addr())); err != nil {
			log.Fatal("could not open embedded server:", err)
		}

		// Wait until the Chromium instance is closed
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, os.Interrupt)
		select {
		case <-sigc:
		case <-ui.Done():
		}
	}
}
