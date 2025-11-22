package main

import (
	"log"
	"machine-marketplace/cmd"
	"machine-marketplace/pkg/database"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Setup graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Shutting down gracefully...")
		if err := database.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
		os.Exit(0)
	}()

	if err := cmd.Main(); err != nil {
		log.Fatal(err)
	}
}
