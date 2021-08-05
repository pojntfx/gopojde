package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"

	"github.com/kataras/compress"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/pojntfx/gopojde/pkg/components"
	"github.com/pojntfx/gopojde/pkg/web"
	"github.com/zserge/lorca"
)

func main() {
	// Browser code
	{
		// Define the routes
		app.Route("/", &components.ManagerHome{})

		// Start the app
		app.RunWhenOnBrowser()
	}

	// Builder code
	{
		if os.Getenv("BUILDER") == "true" {
			// Parse the flags
			build := flag.Bool("build", false, "Create static build")
			out := flag.String("out", "out/gopojde-manager", "Out directory for static build")
			path := flag.String("path", "", "Base path for static build")
			serve := flag.Bool("serve", false, "Build and serve the manager")
			laddr := flag.String("laddr", "localhost:15325", "Address to serve the manager on")

			flag.Parse()

			// Define the handler
			h := &app.Handler{
				Author:          "Felix Pojtinger",
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
			if *build {
				// Deploy under a path
				if *path != "" {
					h.Resources = app.GitHubPages(*path)
				}

				if err := app.GenerateStaticWebsite(*out, h); err != nil {
					log.Fatalf("could not build: %v\n", err)
				}
			}

			// Serve if specified
			if *serve {
				log.Printf("gopojde manager listening on %v\n", *laddr)

				if err := http.ListenAndServe(*laddr, compress.Handler(h)); err != nil {
					log.Fatalf("could not open gopojde manager: %v\n", err)
				}
			}

			return
		}
	}

	// Wrapper code
	{
		// Name the instance
		if runtime.GOOS == "linux" {
			os.Args = append(os.Args, "--class=gopojde Manager")
		}

		// Spawn Chromium instance
		ui, err := lorca.New("", "", 640, 480, os.Args...)
		if err != nil {
			log.Fatal("could not spawn Chromium instance:", err)
		}
		defer ui.Close()

		// Start integrated webserver
		lis, err := net.Listen("tcp", "localhost:0")
		if err != nil {
			log.Fatal("could not listen:", err)
		}
		defer lis.Close()

		// Serve the compiled manager
		root, err := fs.Sub(web.ManagerFS, "manager")
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
