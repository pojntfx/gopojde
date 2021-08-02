package tunnel

import (
	"errors"
	"io"
	"log"
	"net"
	"strings"
	"sync"

	"golang.org/x/crypto/ssh"
)

const (
	tunnelKeySeperator = "->"
)

func marshalTunnelKey(localAddr string, remoteAddr string) string {
	return localAddr + tunnelKeySeperator + remoteAddr
}

func unmarshalTunnelKey(key string) (localAddr string, remoteAddr string, err error) {
	parts := strings.Split(key, tunnelKeySeperator)

	if len(parts) < 2 {
		return "", "", errors.New("could not get key: key does not contain two addresses")
	}

	return parts[0], parts[1], nil
}

type Tunnel struct {
	Key           string
	LocalAddress  string
	RemoteAddress string
}

type SSHConnection struct {
	sshAddr            string
	sshUser            string
	sshAuth            []ssh.AuthMethod
	sshHostKeyCallback func(hostname string, fingerprint string) error
	sshConn            *ssh.Client

	tunnelLock sync.Mutex
	tunnels    map[string]net.Listener
}

func NewSSHConnection(sshAddr string, sshUser string, sshAuth []ssh.AuthMethod, sshHostKeyCallback func(hostname string, fingerprint string) error) *SSHConnection {
	return &SSHConnection{
		sshAddr:            sshAddr,
		sshUser:            sshUser,
		sshAuth:            sshAuth,
		sshHostKeyCallback: sshHostKeyCallback,

		tunnels: map[string]net.Listener{},
	}
}

func (c *SSHConnection) Open() error {
	// Already open; no-op
	if c.sshConn != nil {
		return nil
	}

	conn, err := ssh.Dial("tcp", c.sshAddr, &ssh.ClientConfig{
		User: c.sshUser,
		Auth: c.sshAuth,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return c.sshHostKeyCallback(hostname, ssh.FingerprintSHA256(key))
		},
	})
	if err != nil {
		return err
	}

	c.sshConn = conn

	return nil
}

func (c *SSHConnection) Close() error {
	// Already closed; no-op
	if c.sshConn == nil {
		return nil
	}

	return c.sshConn.Close()
}

func (c *SSHConnection) GetTunnels() ([]Tunnel, error) {
	if c.sshConn == nil {
		return []Tunnel{}, errors.New("could not get tunnels from connection: connection not open")
	}

	tunnels := []Tunnel{}
	for key := range c.tunnels {
		localAddress, remoteAddress, err := unmarshalTunnelKey(key)
		if err != nil {
			return []Tunnel{}, err
		}

		tunnels = append(tunnels, Tunnel{
			Key:           key,
			LocalAddress:  localAddress,
			RemoteAddress: remoteAddress,
		})
	}

	return tunnels, nil
}

func (c *SSHConnection) RemoveTunnel(key string) error {
	c.tunnelLock.Lock()
	defer c.tunnelLock.Unlock()

	if c.sshConn == nil {
		return errors.New("could not close tunnel: connection not open")
	}

	tunnel := c.tunnels[key]
	if tunnel == nil {
		return errors.New("could not close tunnel: tunnel does not exist")
	}

	if err := tunnel.Close(); err != nil {
		return err
	}

	delete(c.tunnels, key)

	return nil
}

func (c *SSHConnection) AddLocalToRemoteTunnel(localAddr string, remoteAddr string) (string, error) {
	c.tunnelLock.Lock()
	defer c.tunnelLock.Unlock()

	if c.sshConn == nil {
		return "", errors.New("could not forward to remote: connection not open")
	}

	// Bind to the desired remote address behind the SSH server
	remote, err := c.sshConn.Listen("tcp", remoteAddr)
	if err != nil {
		return "", err
	}

	// Store the tunnel in memory
	key := marshalTunnelKey(localAddr, remoteAddr)
	c.tunnels[key] = remote

	// Start forwarding
	go func() {
		defer remote.Close()

		removeErroredTunnel := func(err error) {
			log.Printf("closing tunnel %v: %v", key, err)

			if err := c.RemoveTunnel(key); err != nil {
				log.Printf("could not close tunnel %v: %v", key, err)
			}
		}

		// Forward connections to the remote address behind the SSH server to the local address
		for {
			local, err := net.Dial("tcp", localAddr)
			if err != nil {
				removeErroredTunnel(err)

				return
			}

			client, err := remote.Accept()
			if err != nil {
				removeErroredTunnel(err)

				return
			}

			go func(client net.Conn, remote net.Conn) {
				defer client.Close()

				// Handle forwarding errors
				var wg sync.WaitGroup
				wg.Add(2)
				done := make(chan error)

				// Forward from client to remote
				go func(wg *sync.WaitGroup) {
					_, err := io.Copy(client, remote)

					wg.Done()

					done <- err
				}(&wg)

				// Forward from remote to client
				go func(wg *sync.WaitGroup) {
					_, err := io.Copy(remote, client)

					wg.Done()

					done <- err
				}(&wg)

				wg.Wait()
				err := <-done

				// Unexpected close
				if err != nil {
					removeErroredTunnel(err)

					return
				}

				// Nominal close
				if err := c.RemoveTunnel(key); err != nil {
					log.Printf("could not close tunnel %v: %v", key, err)
				}
			}(client, local)
		}
	}()

	return key, nil
}
