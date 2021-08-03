package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"

	api "github.com/pojntfx/gopojde/pkg/api/proto/v1"
	"github.com/pojntfx/gopojde/pkg/tunnel"
	"github.com/pojntfx/gopojde/pkg/web"
	"github.com/zserge/lorca"
	"golang.org/x/crypto/ssh"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ConnectionAndTunnelKey struct {
	ConnectionKey string
	TunnelKey     string
}

func main() {
	if runtime.GOOS == "linux" {
		os.Args = append(os.Args, "--class=gopojde Companion") // No need to quote the `--class` flag, it is already escaped
	}

	ui, err := lorca.New("", "", 480, 640, os.Args...)
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	conn, err := grpc.Dial("localhost:15323", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	sshConns := tunnel.NewSSHConnectionManager()

	client := api.NewInstancesServiceClient(conn)

	if err := ui.Bind("start", func() {
		log.Println("UI is ready!")
	}); err != nil {
		panic(err)
	}

	if err := ui.Bind("getInstances", func() (*api.InstanceSummariesMessage, error) {
		return client.GetInstances(context.Background(), &emptypb.Empty{})
	}); err != nil {
		panic(err)
	}

	if err := ui.Bind("getTunnelsForInstance", func(instanceKey string) ([]tunnel.Tunnel, error) {
		instance, err := sshConns.GetConnection(instanceKey)
		if err != nil {
			return []tunnel.Tunnel{}, err
		}

		return instance.GetTunnels()
	}); err != nil {
		panic(err)
	}

	if err := ui.Bind("forwardToRemote", func(sshPort, sshUser, localAddr, remoteAddr string) (ConnectionAndTunnelKey, error) {
		home, err := os.UserHomeDir()
		if err != nil {
			return ConnectionAndTunnelKey{}, err
		}

		buf, err := ioutil.ReadFile(filepath.Join(home, ".ssh", "id_rsa"))
		if err != nil {
			return ConnectionAndTunnelKey{}, err
		}

		sshKey, err := ssh.ParsePrivateKey(buf)
		if err != nil {
			return ConnectionAndTunnelKey{}, err
		}

		connectionKey, tunnel, err := sshConns.GetOrCreateSSHConnection(net.JoinHostPort("localhost", sshPort), sshUser, []ssh.AuthMethod{ssh.PublicKeys(sshKey)}, func(hostname, fingerprint string) error {
			log.Println(hostname, fingerprint)
			return nil
		})
		if err != nil {
			return ConnectionAndTunnelKey{}, err
		}

		tunnelKey, err := tunnel.AddLocalToRemoteTunnel(localAddr, remoteAddr)
		if err != nil {
			return ConnectionAndTunnelKey{}, err
		}

		return ConnectionAndTunnelKey{
			ConnectionKey: connectionKey,
			TunnelKey:     tunnelKey,
		}, nil
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
