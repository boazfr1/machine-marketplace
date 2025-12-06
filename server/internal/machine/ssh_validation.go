package machine

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

// TriedToConnectForFirstTime validates SSH credentials by attempting a connection
func TriedToConnectForFirstTime(host string, user string, privateKeyString string) error {
	signer, err := ssh.ParsePrivateKey([]byte(privateKeyString))
	if err != nil {
		return fmt.Errorf("failed to parse private key: %v", err)
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", host+":22", config)
	if err != nil {
		return fmt.Errorf("failed to connect: %v", err)
	}
	defer client.Close()

	return nil
}
