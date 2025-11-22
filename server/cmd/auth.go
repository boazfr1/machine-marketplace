package cmd

import (
	"machine-marketplace/internal/user"
	"machine-marketplace/pkg/database"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type AuthService struct {
	Port string `short:"p" long:"port" description:"Port to listen on" default:"3001"`
	database.Config
	AllowedOriginsFrom string `short:"a" long:"allowed-origins-from" description:"Allowed origins from" default:"http://localhost:5173"`
	AllowedOriginsTo   string `short:"t" long:"allowed-origins-to" description:"Allowed origins to" default:"http://localhost:3000"`
}

func (a *AuthService) Setup() user.Module {
	err, queries := database.InitWithConfig(&a.Config)
	if err != nil {
		l.Error("AuthService Execute - failed to initialize database", "error", err)
		return user.Module{}
	}
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{a.AllowedOriginsFrom, a.AllowedOriginsTo},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization"},
		AllowCredentials: true,
	}))
	return user.Module{
		P:  a.Port,
		DB: queries,
		E:  e,
	}
}

func (a *AuthService) Execute(args []string) error {
	l.Info("AuthService Execute - starting auth service", "port", a.Port)

	module := a.Setup()

	module.SetupUserRoutes()

	l.Info("AuthService Execute - starting server", "port", module.P)

	if err := module.E.Start(":" + a.Port); err != nil {
		l.Error("AuthService Execute - failed to start server", "error", err, "port", module.P)
		return err
	}

	return nil
}
