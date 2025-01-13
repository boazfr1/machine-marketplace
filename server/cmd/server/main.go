package main

import (
	"fmt"
	"log"
	"machine-marketplace/internal/routes"
	"machine-marketplace/pkg/database"
	"net/http"
)

const PORT = ":3001"

func main() {

    if err := database.Init(); err != nil {
        log.Fatal("Failed to initialize database:", err)
    }
    defer database.Close()

    mux := http.NewServeMux()
    mux.HandleFunc("/", entryPoint)
    routes.RouteList(mux)

    fmt.Printf("application listening on port %s\n", PORT)
    
    err := http.ListenAndServe(PORT, mux)
    if err != nil {
        log.Fatal(err)
    }

}

func entryPoint(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "welcome")
}