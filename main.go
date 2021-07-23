package main

import (
	"flag"
	"log"

	"github.com/docker/docker/client"
	"github.com/pojntfx/gopojde/pkg/orchestration"
	"github.com/pojntfx/gopojde/pkg/servers"
	"github.com/pojntfx/gopojde/pkg/services"
)

func main() {
	laddr := flag.String("laddr", ":15123", "Listen address")
	wsladdr := flag.String("wsladdr", ":15124", "Listen address (for the WebSocket proxy)")

	flag.Parse()

	docker, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal("could not connect to Docker:", err)
	}

	instancesOrchestration := orchestration.NewInstancesManager(docker)

	instancesService := services.NewInstancesService(instancesOrchestration)

	log.Printf("pojde backend listening on %v (gRPC) and %v (gRPC-Web)\n", *laddr, *wsladdr)

	grpcServer := servers.NewGRPCServer(*laddr, *wsladdr, instancesService)

	log.Fatal(grpcServer.ListenAndServe())
}
