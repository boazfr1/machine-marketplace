package process

import (
	"context"
	"encoding/json"
	db "machine-marketplace/internal/DB/generated"

	kafkaPkg "machine-marketplace/pkg/kafka"

	"github.com/segmentio/kafka-go"
)

type (
	Module struct {
		DB *db.Queries
	}
)

func New(queries *db.Queries) (*Module, error) {
	return &Module{
		DB: queries,
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
	return nil
}