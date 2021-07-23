package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/kataras/compress"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/pojntfx/gopojde/pkg/components"
)

func main() {
	// Client-side code
	{
		// Define the routes
		app.Route("/", &components.Home{})

		// Start the app
		app.RunWhenOnBrowser()
	}

	// Server-/build-side code
	{
		// Parse the flags
		build := flag.Bool("build", false, "Create static build")
		out := flag.String("out", "out/gopojde-frontend", "Out directory for static build")
		path := flag.String("path", "", "Base path for static build")
		serve := flag.Bool("serve", false, "Build and serve the frontend")
		laddr := flag.String("laddr", "localhost:15325", "Address to serve the frontend on")

		flag.Parse()

		// Define the handler
		h := &app.Handler{
			Author:          "Felicitas Pojtinger",
			BackgroundColor: "#151515",
			Description:     "Experimental Go implementation of https://github.com/pojntfx/pojde.",
			Icon: app.Icon{
				Default: "/web/icon.png",
			},
			Keywords: []string{
				// TODO: Add keywords
			},
			LoadingLabel: "Experimental Go implementation of https://github.com/pojntfx/pojde.",
			Name:         "gopojde",
			RawHeaders: []string{
				`<meta property="og:url" content="https://pojntfx.github.io/gopojde/">`,
				`<meta property="og:title" content="gopojde">`,
				`<meta property="og:description" content="Modern network boot server.">`,
				`<meta property="og:image" content="https://pojntfx.github.io/gopojde/web/icon.png">`,
			},
			Styles: []string{
				`https://unpkg.com/@patternfly/patternfly@4.102.2/patternfly.css`,
				`https://unpkg.com/@patternfly/patternfly@4.102.2/patternfly-addons.css`,
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
			log.Printf("gopojde frontend listening on %v\n", *laddr)

			if err := http.ListenAndServe(*laddr, compress.Handler(h)); err != nil {
				log.Fatalf("could not open gopojde frontend: %v\n", err)
			}
		}
	}
}
