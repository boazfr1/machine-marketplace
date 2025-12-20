package machine

import (
	"fmt"
	"log/slog"
	"os"

	"golang.org/x/crypto/ssh"
)

var (
	l = slog.New(slog.NewTextHandler(os.Stdout, nil))
)

// TriedToConnectForFirstTime validates SSH credentials by attempting a connection
func TriedToConnectForFirstTime(host string, user string, privateKeyString string) error {
	l.Info("TriedToConnectForFirstTime - validating SSH connection", "host", host, "user", user)
	
	signer, err := ssh.ParsePrivateKey([]byte(privateKeyString))
	if err != nil {
		l.Error("TriedToConnectForFirstTime - failed to parse private key", "error", err, "host", host)
		return fmt.Errorf("failed to parse private key: %v", err)
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	l.Info("TriedToConnectForFirstTime - attempting SSH connection", "host", host, "user", user)
	client, err := ssh.Dial("tcp", host+":22", config)
	if err != nil {
		l.Error("TriedToConnectForFirstTime - failed to connect", "error", err, "host", host, "user", user)
		return fmt.Errorf("failed to connect: %v", err)
	}
	defer client.Close()

	l.Info("TriedToConnectForFirstTime - SSH connection successful", "host", host, "user", user)
	return nil
}
