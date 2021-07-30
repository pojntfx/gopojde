package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"

	"github.com/pojntfx/gopojde/pkg/web"
	"github.com/zserge/lorca"
)

func main() {
	if runtime.GOOS == "linux" {
		os.Args = append(os.Args, "--class=gopojde Companion") // No need to quote the `--class` flag, it is already escaped
	}

	ui, err := lorca.New("", "", 640, 480, os.Args...)
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	if err := ui.Bind("start", func() {
		log.Println("UI is ready!")
	}); err != nil {
		panic(err)
	}

	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}
	defer lis.Close()

	go http.Serve(lis, http.FileServer(http.FS(web.CompanionFS)))

	if err := ui.Load(fmt.Sprintf("http://%v/companion", lis.Addr())); err != nil {
		panic(err)
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)
	select {
	case <-sigc:
	case <-ui.Done():
	}
}
