package tunnel

import (
	"errors"
	"io"
	"log"
	"net"
	"sync"

	"golang.org/x/crypto/ssh"
)

type SSHTunnelManager struct {
	sshAddr            string
	sshUser            string
	sshAuth            []ssh.AuthMethod
	sshHostKeyCallback func(hostname string, fingerprint string) error
	sshConn            *ssh.Client
}

func NewSSHTunnelManager(sshAddr string, sshUser string, sshAuth []ssh.AuthMethod, sshHostKeyCallback func(hostname string, fingerprint string) error) *SSHTunnelManager {
	return &SSHTunnelManager{
		sshAddr:            sshAddr,
		sshUser:            sshUser,
		sshAuth:            sshAuth,
		sshHostKeyCallback: sshHostKeyCallback,
	}
}

func (m *SSHTunnelManager) Open() error {
	// Already open; no-op
	if m.sshConn != nil {
		return nil
	}

	conn, err := ssh.Dial("tcp", m.sshAddr, &ssh.ClientConfig{
		User: m.sshUser,
		Auth: m.sshAuth,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return m.sshHostKeyCallback(hostname, ssh.FingerprintSHA256(key))
		},
	})
	if err != nil {
		return err
	}

	m.sshConn = conn

	return nil
}

func (m *SSHTunnelManager) Close() error {
	// Already closed; no-op
	if m.sshConn == nil {
		return nil
	}

	return m.sshConn.Close()
}

func (m *SSHTunnelManager) ForwardToRemote(localAddr string, remoteAddr string) error {
	if m.sshConn == nil {
		return errors.New("could not forward to remote: connection not open")
	}

	// Bind to the desired remote address behind the SSH server
	remote, err := m.sshConn.Listen("tcp", remoteAddr)
	if err != nil {
		return err
	}
	defer remote.Close()

	// Forward connections to the remote address behind the SSH server to the local address
	// TODO: Store opened connections in memory
	for {
		local, err := net.Dial("tcp", localAddr)
		if err != nil {
			return err
		}

		client, err := remote.Accept()
		if err != nil {
			return err
		}

		go func(client net.Conn, remote net.Conn) {
			defer client.Close()

			var wg sync.WaitGroup
			wg.Add(2)

			go func(wg *sync.WaitGroup) {
				if _, err := io.Copy(client, remote); err != nil {
					log.Println("could not forward from client to remote")
				}

				wg.Done()
			}(&wg)

			go func(wg *sync.WaitGroup) {
				if _, err := io.Copy(remote, client); err != nil {
					log.Println("could not forward from remote to client")
				}

				wg.Done()
			}(&wg)

			wg.Wait()
		}(client, local)
	}
}
