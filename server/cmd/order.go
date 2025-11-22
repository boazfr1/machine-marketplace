package cmd

import (
	"machine-marketplace/pkg/database"
	"machine-marketplace/pkg/order"
)

type OrderService struct {
	Port string `short:"p" long:"port" description:"Port to listen on" default:"3002"`
	database.Config
}

func (o *OrderService) Execute(args []string) error {
	l.Info("OrderService Execute - starting order service", "port", o.Port)

	// Initialize database with configuration
	if err := database.InitWithConfig(&o.Config); err != nil {
		l.Error("OrderService Execute - failed to initialize database", "error", err)
		return err
	}

	err := order.New(o.Port)
	if err != nil {
		l.Error("OrderService Execute - failed to start order service", "error", err, "port", o.Port)
		return err
	}

	l.Info("OrderService Execute - order service started successfully", "port", o.Port)
	return nil
}
