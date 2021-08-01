package tunnel

import (
	"io"
	"log"
	"net"
	"sync"

	"golang.org/x/crypto/ssh"
)

func ForwardToRemoteSSH(localAddr string, remoteAddr string, sshAddr string, sshUser string, sshPrivateKey ssh.Signer) error {
	// Connect to the SSH server
	conn, err := ssh.Dial("tcp", sshAddr, &ssh.ClientConfig{
		User: sshUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(sshPrivateKey),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // TODO: Use `known_hosts` and interactive validation
	})
	if err != nil {
		return err
	}
	defer conn.Close()

	// Bind to the desired remote address behind the SSH server
	remote, err := conn.Listen("tcp", remoteAddr)
	if err != nil {
		return err
	}
	defer remote.Close()

	// Forward connections to the remote address behind the SSH server to the local address
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
