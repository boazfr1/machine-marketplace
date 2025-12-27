package process

import (
	db "machine-marketplace/internal/DB/generated"
	"machine-marketplace/pkg/process"
	"machine-marketplace/pkg/kafka"
	"os"

	"log/slog"

	"context"
)

type (
	Module struct {
		P             string
		DB            *db.Queries
		KafkaConsumer *kafka.Consumer
		Process       *process.Module
		ctx           context.Context
	}
)

var (
	l = slog.New(slog.NewTextHandler(os.Stdout, nil))
)

func New(port string, queries *db.Queries, kafkaAddress string) (Module, error) {

	kafkaConsumer, err := kafka.NewConsumer(kafkaAddress, "process-consumer-group")
	if err != nil {
		l.Error("ProcessService Execute - failed to initialize Kafka consumer", "error", err)
		return Module{}, err
	}

	process, err := process.New(queries)
	if err != nil {
		l.Error("ProcessService Execute - failed to initialize process", "error", err)
		return Module{}, err
	}
	defer func() {
		if err := kafkaConsumer.Close(); err != nil {
			l.Error("ProcessService Execute - failed to close Kafka consumer", "error", err)
		}
	}()

	return Module{
		P:             port,
		DB:            queries,
		KafkaConsumer: kafkaConsumer,
		Process:       process,
		ctx:           context.Background(),
	}, nil
}

func (m *Module) ReadAndProcess() error {
	for {
		message, err := m.KafkaConsumer.ReadMessage(m.ctx)
		if err != nil {
			l.Error("ReadAndProcess - failed to read message", "error", err)
			return err
		}

		m.Process.ProcessMessage(message)
		if err != nil {
			l.Error("ReadAndProcess - failed to process message", "error", err)
			return err
		}

		l.Info("ReadAndProcess - successfully processed message", "message", string(message.Value))
	}
}
