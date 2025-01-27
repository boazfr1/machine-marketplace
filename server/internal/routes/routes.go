package routes

import (
	"encoding/json"
	"net/http"

	machine "machine-marketplace/internal/machine"
	middleware "machine-marketplace/internal/middleware"
	user "machine-marketplace/internal/user"
)

func RouteList(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/health", middleware.Get(health))
	mux.HandleFunc("/api/v1/sign-up", middleware.Post(user.SignUp))
	mux.HandleFunc("/api/v1/login", middleware.Post(user.Login))
	mux.HandleFunc("/api/v1/logout", middleware.Post(user.Logout))
	mux.HandleFunc("/api/v1/user", middleware.Get(user.User))
	mux.HandleFunc("/api/v1/machine", middleware.Get(machine.ListOfFreeMachines))
	mux.HandleFunc("/api/v1/machine/create", middleware.Post(machine.CreateMachine))
	mux.HandleFunc("/api/v1/machine/connect", middleware.Post(machine.WebSocketHandler))
	mux.HandleFunc("/api/v1/machine/my-machines", middleware.Get(machine.GetMyMachines))

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
