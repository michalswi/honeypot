package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
)

var sshPort = getEnv("SSH_PORT", "2222")
var keyName = "private"

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func handleSSHConnection(conn net.Conn, config *ssh.ServerConfig) {
	defer conn.Close()

	sshConn, chans, reqs, err := ssh.NewServerConn(conn, config)
	if err != nil {
		log.Printf("Error during SSH handshake: %v\n", err)
		return
	}

	log.Printf("SSH login attempt [remote addr] from %s, username: %s, password: %s, client version: %s",
		conn.RemoteAddr(),
		sshConn.User(),
		sshConn.Permissions.Extensions["password"],
		sshConn.ClientVersion(),
	)

	go ssh.DiscardRequests(reqs)

	for newChannel := range chans {
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}

		channel, _, err := newChannel.Accept()
		if err != nil {
			log.Printf("Error accepting channel: %v\n", err)
			continue
		}
		channel.Close()
	}
}

func main() {
	config := &ssh.ServerConfig{
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			permissions := &ssh.Permissions{
				Extensions: map[string]string{
					"password": string(pass),
				},
			}
			return permissions, nil
		},
	}

	privateBytes, err := ioutil.ReadFile(keyName)
	if err != nil {
		log.Fatal("Failed to load private key:", err)
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Fatal("Failed to parse private key:", err)
	}

	config.AddHostKey(private)

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", sshPort))
	if err != nil {
		log.Fatal("Error listening for connection:", err)
	}

	defer listener.Close()

	log.Printf("SSH server started. Listening on port %s", sshPort)

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}

		go handleSSHConnection(conn, config)
	}
}
