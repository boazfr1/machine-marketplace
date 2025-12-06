package cmd

import (
	"machine-marketplace/internal/order"
	"machine-marketplace/pkg/database"
	"machine-marketplace/pkg/kafka"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type OrderService struct {
	Port string `short:"p" long:"port" description:"Port to listen on" default:"3002"`
	database.Config
	AllowedOriginsFrom string `short:"a" long:"allowed-origins-from" description:"Allowed origins from" default:"http://localhost:5173"`
	AllowedOriginsTo   string `short:"t" long:"allowed-origins-to" description:"Allowed origins to" default:"http://localhost:3000"`
}

func (o *OrderService) Setup() order.Module {

	err, queries := database.InitWithConfig(&o.Config)
	if err != nil {
		l.Error("OrderService Execute - failed to initialize database", "error", err)
		return order.Module{}
	}

	// Initialize Kafka producer
	kafkaProducer, err := kafka.New()
	if err != nil {
		l.Error("OrderService Execute - failed to initialize Kafka producer", "error", err)
		return order.Module{}
	}

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{o.AllowedOriginsFrom, o.AllowedOriginsTo},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization"},
		AllowCredentials: true,
	}))
	return order.Module{
		P:             o.Port,
		DB:            queries,
		E:             e,
		KafkaProducer: kafkaProducer,
	}
}

func (o *OrderService) Execute(args []string) error {
	l.Info("OrderService Execute - starting order service", "port", o.Port)

	module := o.Setup()

	module.SetupOrderRoutes()

	l.Info("OrderService Execute - starting server", "port", module.P)

	if err := module.E.Start(":" + o.Port); err != nil {
		l.Error("OrderService Execute - failed to start server", "error", err, "port", module.P)
		return err
	}

	l.Info("OrderService Execute - order service started successfully", "port", o.Port)
	return nil
}
