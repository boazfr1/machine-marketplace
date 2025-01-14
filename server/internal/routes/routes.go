package routes

import (
	"encoding/json"
	"net/http"

	"machine-marketplace/internal/middleware"
	"machine-marketplace/internal/user"
)

func RouteList(mux *http.ServeMux) {
    mux.HandleFunc("/api/v1/health", middleware.Get(health));
	mux.HandleFunc("/api/v1/boaz", middleware.Get(boaz));
	mux.HandleFunc("/api/v1/sign-up", middleware.Post(user.SignUp))
	mux.HandleFunc("/api/v1/login", middleware.Post(user.Login))

}

func health(res http.ResponseWriter, req *http.Request) {

    data := map[string]string{
		"massage": "be kind",
	};

	js, err := json.Marshal(data);
	if err != nil {
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(js);

}

func boaz(res http.ResponseWriter, req *http.Request) {
	
}