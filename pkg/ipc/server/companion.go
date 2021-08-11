package server

import (
	"context"
	"errors"
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

type CompanionIPCServer struct {
	daemon  api.InstancesServiceClient
	address string

	sshConnectionManager     *tunnel.SSHConnectionManager
	sshConnectionManagerLock sync.Mutex
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

	if err := ui.Bind(shared.CreateSSHConnection, c.CreateSSHConnection); err != nil {
		return err
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

func (c *CompanionIPCServer) GetInstances() ([]shared.Instance, error) {
	if c.daemon == nil {
		return []shared.Instance{}, errors.New("could not get instances: not connected to daemon")
	}

	// Get all instances
	instances, err := c.daemon.GetInstances(context.Background(), &emptypb.Empty{})
	if err != nil {
		return []shared.Instance{}, err
	}

	// Reduce instances to relevant options
	res := []shared.Instance{}
	for _, instance := range instances.GetInstances() {
		res = append(res, shared.Instance{
			ID:   instance.GetInstanceID().GetName(),
			Name: instance.GetInstanceID().GetName(),
		})
	}

	return res, nil
}

func (c *CompanionIPCServer) CreateSSHConnection(
	instanceID string,
	privateKey string,
	passwordGetterFunc func() string,
	hostKeyValidatorFunc func(hostname, fingerprint string) error,
) (string, error) {
	if c.daemon == nil {
		return "", errors.New("could not get instances: not connected to daemon")
	}

	// Get all instances
	instances, err := c.daemon.GetInstances(context.Background(), &emptypb.Empty{})
	if err != nil {
		return "", err
	}

	// Get the relevant instance
	targetPort := 0
	targetUser := ""
	found := false
	for _, instance := range instances.GetInstances() {
		if instance.GetInstanceID().GetName() == instanceID {
			ports := instance.GetPorts()

			for _, port := range ports {
				if port.GetService() == "ssh" {
					targetPort = int(port.GetPort())

					config, err := c.daemon.GetInstanceConfig(context.Background(), instance.GetInstanceID())
					if err != nil {
						return "", err
					}

					targetUser = config.GetUserName()

					found = true

					break
				}
			}
		}
	}

	if !found {
		return "", errors.New("could not find SSH credentials for instance")
	}

	// Get hostname from Docker
	u, err := url.Parse(c.address)
	if err != nil {
		return "", err
	}

	// Parse the SSH key
	sshKey, err := ssh.ParsePrivateKey([]byte(privateKey))
	if err != nil {
		if err.Error() == (&ssh.PassphraseMissingError{}).Error() {
			sshKey, err = ssh.ParsePrivateKeyWithPassphrase([]byte(privateKey), []byte(passwordGetterFunc()))
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}

	// Create the SSH connection manager if it doesn't already exist
	c.sshConnectionManagerLock.Lock()
	defer c.sshConnectionManagerLock.Unlock()

	if c.sshConnectionManager == nil {
		c.sshConnectionManager = tunnel.NewSSHConnectionManager()
	}

	// Create the SSH connection
	key, _, err := c.sshConnectionManager.GetOrCreateSSHConnection(
		net.JoinHostPort(u.Host, strconv.Itoa(targetPort)),
		targetUser,
		[]ssh.AuthMethod{ssh.PublicKeys(sshKey)},
		func(hostname, fingerprint string) error {
			return hostKeyValidatorFunc(hostname, fingerprint)
		})
	if err != nil {
		return "", err
	}

	return key, nil
}
