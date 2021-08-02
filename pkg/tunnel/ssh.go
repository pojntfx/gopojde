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
	tunnelKeySeparator     = "->"
	sshConnectionSeparator = "@"
)

func marshalTunnelKey(localAddr string, remoteAddr string) string {
	return localAddr + tunnelKeySeparator + remoteAddr
}

func unmarshalTunnelKey(key string) (localAddr string, remoteAddr string, err error) {
	parts := strings.Split(key, tunnelKeySeparator)

	if len(parts) < 2 {
		return "", "", errors.New("could not get key: key does not contain two addresses")
	}

	return parts[0], parts[1], nil
}

func marshalSSHConnectionKey(addr string, user string) string {
	return user + sshConnectionSeparator + addr
}

func unmarshalSSHConnectionKey(key string) (addr string, user string, err error) {
	parts := strings.Split(key, sshConnectionSeparator)

	if len(parts) < 2 {
		return "", "", errors.New("could not get key: key does not contain user and port")
	}

	return parts[0], parts[1], nil
}

type Connection struct {
	Key     string
	User    string
	Address string
}

type Tunnel struct {
	Key           string
	LocalAddress  string
	RemoteAddress string
}

type SSHConnectionManager struct {
	connectionLock sync.Mutex
	connections    map[string]*SSHConnection
}

func NewSSHConnectionManager() *SSHConnectionManager {
	return &SSHConnectionManager{
		connections: map[string]*SSHConnection{},
	}
}

func (m *SSHConnectionManager) GetOrCreateSSHConnection(sshAddr string, sshUser string, sshAuth []ssh.AuthMethod, sshHostKeyCallback func(hostname string, fingerprint string) error) (string, *SSHConnection, error) {
	m.connectionLock.Lock()
	defer m.connectionLock.Unlock()

	// Store the connection in memory if it does not already exist
	key := marshalSSHConnectionKey(sshAddr, sshUser)
	conn := m.connections[key]
	if conn == nil {
		// Create a new SSH connection
		conn = NewSSHConnection(sshAddr, sshUser, sshAuth, sshHostKeyCallback)

		// Open the SSH connection
		if err := conn.Open(); err != nil {
			return "", nil, err
		}

		m.connections[key] = conn
	}

	return key, conn, nil
}

func (m *SSHConnectionManager) GetConnections() ([]Connection, error) {
	connections := []Connection{}
	for key := range m.connections {
		addr, user, err := unmarshalSSHConnectionKey(key)
		if err != nil {
			return []Connection{}, err
		}

		connections = append(connections, Connection{
			Key:     key,
			User:    user,
			Address: addr,
		})
	}

	return connections, nil
}

func (m *SSHConnectionManager) GetConnection(key string) (*SSHConnection, error) {
	connection := m.connections[key]
	if connection == nil {
		return &SSHConnection{}, errors.New("could not get connection: connection does not exist")
	}

	return connection, nil
}

func (c *SSHConnectionManager) RemoveConnection(key string) error {
	c.connectionLock.Lock()
	defer c.connectionLock.Unlock()

	connections := c.connections[key]
	if connections == nil {
		return errors.New("could not remove connection: connection does not exist")
	}

	if err := connections.Close(); err != nil {
		return err
	}

	delete(c.connections, key)

	return nil
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

	for _, conn := range c.tunnels {
		_ = conn.Close() // Ignore closing errors; we're tearing down
	}

	if err := c.sshConn.Close(); err != nil {
		return nil
	}

	c.sshConn = nil

	return nil
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
		return errors.New("could not remove tunnel: connection not open")
	}

	tunnel := c.tunnels[key]
	if tunnel == nil {
		return errors.New("could not remove tunnel: tunnel does not exist")
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

	// Store the tunnel in memory if it does not already exist
	key := marshalTunnelKey(localAddr, remoteAddr)
	remote := c.tunnels[key]
	if remote == nil {
		// Bind to the desired remote address behind the SSH server
		r, err := c.sshConn.Listen("tcp", remoteAddr)
		if err != nil {
			return "", err
		}

		remote = r
		c.tunnels[key] = remote
	} else {
		// Tunnel is already forwarded
		return key, nil
	}

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
