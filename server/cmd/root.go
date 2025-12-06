package cmd

import (
	"log/slog"
	"os"

	"github.com/jessevdk/go-flags"
)

type (
	rootCmd struct {
		OrderService   OrderService   `command:"order-service" description:"Order service"`
		ProcessService ProcessService `command:"process-service" description:"Process service"`
		AuthService    AuthService    `command:"auth-service" description:"Auth service"`
	}
)

var (
	l = slog.New(slog.NewTextHandler(os.Stdout, nil))
)

func Main() error {
	var cmd rootCmd
	parser := flags.NewParser(&cmd, flags.Default)

	_, err := parser.Parse()
	if err != nil {
		l.Error("Main - failed to parse flags", "error", err)
		return err
	}

	l.Info("Main - application started successfully")
	return nil
}
