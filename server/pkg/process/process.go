package process

import (
	"context"
	"encoding/json"
	db "machine-marketplace/internal/DB/generated"

	kafkaPkg "machine-marketplace/pkg/kafka"
	sshconection "machine-marketplace/pkg/sshConection"

	"github.com/segmentio/kafka-go"
)

type (
	Module struct {
		DB        *db.Queries
		OpenaiKey string
		Model     string
	}

	Config struct {
		DB        *db.Queries
		OpenaiKey string
		Model     string
	}
)

func New(config Config) (*Module, error) {
	return &Module{
		DB:        config.DB,
		OpenaiKey: config.OpenaiKey,
		Model:     config.Model,
	}, nil
}

func (m *Module) ProcessMessage(message kafka.Message) error {
	var event kafkaPkg.PurchaseEvent
	if err := json.Unmarshal(message.Value, &event); err != nil {
		return err
	}

	machineDetails, err := m.DB.GetMachineByID(context.Background(), event.MachineID)
	if err != nil {
		return err
	}

	err = m.setUpVirtualMachine(machineDetails)
	if err != nil {
		return err
	}

	return nil
}

func (m *Module) setUpVirtualMachine(machineDetails db.Machine) error {
	ssh, err := sshconection.New(sshconection.Config{
		Host:      machineDetails.Host,
		User:      machineDetails.SshUser,
		Key:       machineDetails.Key.String,
		OpenaiKey: m.OpenaiKey,
		Model:     m.Model,
	})
	if err != nil {
		return err
	}
	defer ssh.Close()

	result, err := ssh.RunSSHCommand("echo 'Hello, World!'")
	if err != nil {
		return err
	}

	_ = result // Use result if needed for logging or further processing

	return nil
}
