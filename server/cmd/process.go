package cmd

import (
	"machine-marketplace/pkg/database"
	"machine-marketplace/pkg/process"
)

type ProcessService struct {
	Port string `short:"p" long:"port" description:"Port to listen on" default:"3003"`
	database.Config
}

func (p *ProcessService) Execute(args []string) error {
	l.Info("ProcessService Execute - starting process service", "port", p.Port)

	if err := database.InitWithConfig(&p.Config); err != nil {
		l.Error("ProcessService Execute - failed to initialize database", "error", err)
		return err
	}

	err := process.New(p.Port)
	if err != nil {
		l.Error("ProcessService Execute - failed to start process service", "error", err, "port", p.Port)
		return err
	}

	l.Info("ProcessService Execute - process service started successfully", "port", p.Port)
	return nil
}
