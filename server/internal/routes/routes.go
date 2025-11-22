package routes

import (
	"encoding/json"
	"net/http"

	machine "machine-marketplace/internal/machine"
	middleware "machine-marketplace/internal/middleware"
)

func RouteList(mux *http.ServeMux) {

	mux.HandleFunc("/api/v1/machine", middleware.GetWithAuth(machine.ListOfFreeMachines))
	mux.HandleFunc("/api/v1/machine/create", middleware.PostWithAuth(machine.CreateMachine))
	mux.HandleFunc("/api/v1/machine/connect", middleware.GetWithAuth(machine.WebSocketHandler))
	mux.HandleFunc("/api/v1/machine/owned-machines", middleware.GetWithAuth(machine.GetOwnedMachines))
	mux.HandleFunc("/api/v1/machine/bought-machines", middleware.GetWithAuth(machine.GetBoughtMachines))
}

func health(res http.ResponseWriter, req *http.Request) {

	data := map[string]string{
		"massage": "be kind",
	}

	js, err := json.Marshal(data)
	if err != nil {
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(js)

}
