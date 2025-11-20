package cmd

import (
	"fmt"
	"log"

	"machine-marketplace/internal/middleware"
	"machine-marketplace/internal/routes"
	"machine-marketplace/pkg/database"
	"net/http"
	"github.com/jessevdk/go-flags"
)

type (
	rootFlags struct {
		Port string `short:"p" long:"port" description:"Port to listen on" default:"3001"`
	}

	rootCmd struct {
		OrderService OrderService `command:"order-service" description:"Order service"`
		ProcessService ProcessService `command:"process-service" description:"Process service"`
	}

) 

const PORT = ":3001"

func Main() error {

	if err := database.Init(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer database.Close()

	if err := database.SetupDatabase(); err != nil {
		log.Fatal("Failed to setup database:", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", entryPoint)
	routes.RouteList(mux)

	muxWithCORS := middleware.EnableCORS(mux)

	fmt.Printf("application listening on port %s\n", PORT)

	err := http.ListenAndServe(PORT, muxWithCORS)
	if err != nil {
		log.Fatal(err)
	}

}

func entryPoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "welcome")
}
