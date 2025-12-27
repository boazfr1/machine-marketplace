package cmd

import (
	"machine-marketplace/internal/process"
	"machine-marketplace/pkg/database"
)

type ProcessService struct {
	Port         string `short:"p" long:"port" description:"Port to listen on" default:"3003"`
	KafkaAddress string `short:"k" long:"kafka-address" description:"Kafka address" default:"localhost:9092"`
	database.Config
}

func (p *ProcessService) Execute(args []string) error {
	l.Info("ProcessService Execute - starting process service", "port", p.Port)

	err, queries := database.InitWithConfig(&p.Config)
	if err != nil {
		l.Error("OrderService Execute - failed to initialize database", "error", err)
		return err
	}

	module, err := process.New(p.Port, queries, p.KafkaAddress)
	if err != nil {
		l.Error("ProcessService Execute - failed to start process service", "error", err, "port", p.Port)
		return err
	}

	defer func() {
		if err := module.KafkaConsumer.Close(); err != nil {
			l.Error("ProcessService Execute - failed to close Kafka consumer", "error", err)
		}
	}()

	l.Info("ProcessService Execute - process service started successfully", "port", p.Port)
	return nil
}
