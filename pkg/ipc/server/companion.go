package server

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"net"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/pojntfx/go-app-grpc-chat-frontend-web/pkg/websocketproxy"
	api "github.com/pojntfx/gopojde/pkg/api/proto/v1"
	"github.com/pojntfx/gopojde/pkg/ipc/shared"
	"github.com/pojntfx/gopojde/pkg/tunnel"
	"github.com/zserge/lorca"
	"golang.org/x/crypto/ssh"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	targetUser = "root"
)

type CompanionIPCServer struct {
	daemon  api.InstancesServiceClient
	address string

	sshConnectionManager     *tunnel.SSHConnectionManager
	sshConnectionManagerLock sync.Mutex

	passwordGetterFunc   func() string
	hostKeyValidatorFunc func(hostname, fingerprint string) error
}

func NewCompanionIPC() *CompanionIPCServer {
	return &CompanionIPCServer{}
}

func (c *CompanionIPCServer) Bind(ui lorca.UI) error {
	if err := ui.Bind(shared.OpenKey, c.Open); err != nil {
		return err
	}

	if err := ui.Bind(shared.GetInstancesKey, c.GetInstances); err != nil {
		return err
	}

	if err := ui.Bind(shared.ForwardFromLocalToRemoteKey, c.ForwardFromLocalToRemote); err != nil {
		return err
	}

	c.passwordGetterFunc = func() string {
		return ui.Eval(shared.PasswordGetterKey + "()").String()
	}

	c.hostKeyValidatorFunc = func(hostname, fingerprint string) error {
		return ui.Eval(shared.HostKeyValidatorKey + fmt.Sprintf(`("%v", "%v")`, template.JSEscapeString(hostname), template.JSEscapeString(fingerprint))).Err()
	}

	return nil
}

func (c *CompanionIPCServer) Open(address string) error {
	conn, err := grpc.Dial(address, grpc.WithContextDialer(websocketproxy.NewWebSocketProxyClient(time.Minute).Dialer), grpc.WithInsecure())
	if err != nil {
		return err
	}

	c.daemon = api.NewInstancesServiceClient(conn)
	c.address = address

	return nil
}

func (c *CompanionIPCServer) GetInstances(privateKey string) ([]shared.Instance, error) {
	if c.daemon == nil {
		return []shared.Instance{}, errors.New("could not get instances: not connected to daemon")
	}

	// Get all instances
	instances, err := c.daemon.GetInstances(context.Background(), &emptypb.Empty{})
	if err != nil {
		return []shared.Instance{}, err
	}

	// Get hostname from Docker
	u, err := url.Parse(c.address)
	if err != nil {
		return []shared.Instance{}, err
	}

	// Parse the SSH key
	sshKey, err := ssh.ParsePrivateKey([]byte(privateKey))
	if err != nil {
		if err.Error() == (&ssh.PassphraseMissingError{}).Error() {
			sshKey, err = ssh.ParsePrivateKeyWithPassphrase([]byte(privateKey), []byte(c.passwordGetterFunc()))
			if err != nil {
				return []shared.Instance{}, err
			}
		} else {
			return []shared.Instance{}, err
		}
	}

	// Create the SSH connection manager if it doesn't already exist
	c.sshConnectionManagerLock.Lock()
	defer c.sshConnectionManagerLock.Unlock()

	if c.sshConnectionManager == nil {
		c.sshConnectionManager = tunnel.NewSSHConnectionManager()
	}

	// Reduce instances to relevant options
	res := []shared.Instance{}
	for _, instance := range instances.GetInstances() {
		// Get the SSH port for the instance
		ports := instance.GetPorts()
		sshPort := -1
		for _, port := range ports {
			if port.GetService() == "ssh" {
				sshPort = int(port.GetPort())

				break
			}
		}

		if sshPort == -1 {
			return []shared.Instance{}, errors.New("could not find SSH connection details for instance")
		}

		// Create the SSH connection
		_, conn, err := c.sshConnectionManager.GetOrCreateSSHConnection(
			net.JoinHostPort(u.Hostname(), strconv.Itoa(sshPort)),
			// TODO: `pojde` instances currently only configure SSH access for the root user,
			// but in the future the non-root user could be queried with config.GetUserName()
			targetUser,
			[]ssh.AuthMethod{ssh.PublicKeys(sshKey)},
			func(hostname, fingerprint string) error {
				return c.hostKeyValidatorFunc(hostname, fingerprint)
			})
		if err != nil {
			return []shared.Instance{}, err
		}

		// Query the tunnels for the instance
		rawTunnels, err := conn.GetTunnels()
		if err != nil {
			return []shared.Instance{}, err
		}

		// Reduce tunnels to relevant options
		tunnels := []shared.Tunnel{}
		for _, tunnel := range rawTunnels {
			tunnels = append(tunnels, shared.Tunnel{
				ID:            tunnel.Key,
				LocalAddress:  tunnel.LocalAddress,
				RemoteAddress: tunnel.RemoteAddress,
			})
		}

		res = append(res, shared.Instance{
			ID:      instance.GetInstanceID().GetName(),
			Name:    instance.GetInstanceID().GetName(),
			Tunnels: tunnels,
		})
	}

	return res, nil
}

func (c *CompanionIPCServer) ForwardFromLocalToRemote(instanceID string, localAddr string, remoteAddr string) (string, error) {
	if c.daemon == nil {
		return "", errors.New("could not forward port: not connected to daemon")
	}

	if c.sshConnectionManager == nil {
		return "", errors.New("could not forward port: not connected to instance")
	}

	// FIXME: Add `ConnectionID` (i.e. root@localhost:5005) to the instances struct and use that for a fully-qualified ID
	conn, err := c.sshConnectionManager.GetConnection(instanceID)
	if err != nil {
		return "", err
	}

	return conn.AddLocalToRemoteTunnel(localAddr, remoteAddr)
}
