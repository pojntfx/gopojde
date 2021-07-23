package servers

import (
	"net"
	"net/http"
	"sync"

	"github.com/pojntfx/bofied/pkg/websocketproxy"
	api "github.com/pojntfx/gopojde/pkg/api/proto/v1"
	"github.com/pojntfx/gopojde/pkg/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	listenAddress          string
	websocketListenAddress string

	instancesService *services.InstancesService
}

func NewGRPCServer(listenAddress string, websocketListenAddress string, instancesService *services.InstancesService) *GRPCServer {
	return &GRPCServer{
		listenAddress:          listenAddress,
		websocketListenAddress: websocketListenAddress,

		instancesService: instancesService,
	}
}

func (s *GRPCServer) ListenAndServe() error {
	proxy := websocketproxy.NewWebSocketProxyServer()
	listener, err := net.Listen("tcp", s.listenAddress)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	reflection.Register(server)

	api.RegisterInstancesServiceServer(server, s.instancesService)

	doneChan := make(chan struct{})
	errChan := make(chan error)

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		wg.Wait()

		close(doneChan)
	}()

	go func() {
		if err := server.Serve(listener); err != nil {
			errChan <- err
		}

		wg.Done()
	}()

	go func() {
		if err := server.Serve(proxy); err != nil {
			errChan <- err
		}

		wg.Done()
	}()

	go func() {
		if err := http.ListenAndServe(s.websocketListenAddress, proxy); err != nil {
			errChan <- err
		}

		wg.Done()
	}()

	select {
	case <-doneChan:
		return nil
	case <-errChan:
		return err
	}
}
