package sshconection

import (
	"context"
	"fmt"
	llmPkg "machine-marketplace/pkg/llm"

	"golang.org/x/crypto/ssh"
)

const maxRetries = 3

type (
	Module struct {
		Client    *ssh.Client
		OpenaiKey string
		Model     string
	}

	Config struct {
		Host      string
		User      string
		Key       string
		OpenaiKey string
		Model     string
	}

	CommandResult struct {
		Output  string
		Success bool
	}
)

func New(config Config) (*Module, error) {
	signer, err := convertStringToSigner(config.Key)
	if err != nil {
		return nil, err
	}
	sshConfig := &ssh.ClientConfig{
		User: config.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", config.Host+":22", sshConfig)
	if err != nil {
		return nil, err
	}

	return &Module{
		Client:    client,
		OpenaiKey: config.OpenaiKey,
		Model:     config.Model,
	}, nil
}

func (m *Module) Close() error {
	if m.Client != nil {
		return m.Client.Close()
	}
	return nil
}

func (m *Module) RunSSHCommand(command string) (*CommandResult, error) {
	return m.runSSHCommandWithRetry(context.Background(), command, 0)
}

func (m *Module) RunSSHCommandWithContext(ctx context.Context, command string) (*CommandResult, error) {
	return m.runSSHCommandWithRetry(ctx, command, 0)
}

func (m *Module) runSSHCommandWithRetry(ctx context.Context, command string, attempt int) (*CommandResult, error) {
	session, err := m.Client.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	output, err := session.CombinedOutput(command)
	if err == nil {
		return &CommandResult{
			Output:  string(output),
			Success: true,
		}, nil
	}

	// Command failed, try to fix with LLM
	if attempt >= maxRetries {
		return &CommandResult{
			Output:  string(output),
			Success: false,
		}, fmt.Errorf("max retries exceeded: %w", err)
	}

	fixResponse, llmErr := llmPkg.HandleErrorWithLLM(ctx, m.OpenaiKey, command, string(output), err, m.Model)
	if llmErr != nil {
		return &CommandResult{
			Output:  string(output),
			Success: false,
		}, err
	}

	if !fixResponse.CanFix {
		return &CommandResult{
			Output:  string(output),
			Success: false,
		}, fmt.Errorf("unfixable error: %s - %w", fixResponse.Explanation, err)
	}

	// Retry with the new command suggested by LLM
	return m.runSSHCommandWithRetry(ctx, fixResponse.NewCommand, attempt+1)
}

func convertStringToSigner(key string) (ssh.Signer, error) {
	signer, err := ssh.ParsePrivateKey([]byte(key))
	if err != nil {
		return nil, err
	}
	return signer, nil
}
