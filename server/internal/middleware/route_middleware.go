package middleware

import (
	"log/slog"
	"net/http"
	"os"
)

var (
	routeLogger = slog.New(slog.NewTextHandler(os.Stdout, nil))
)

func Get(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			routeLogger.Error("Get - method not allowed", "method", req.Method, "path", req.URL.Path)
			http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		routeLogger.Info("Get - processing GET request", "path", req.URL.Path)
		handler(res, req)
	}
}

func Post(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			routeLogger.Error("Post - method not allowed", "method", req.Method, "path", req.URL.Path)
			http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		routeLogger.Info("Post - processing POST request", "path", req.URL.Path)
		handler(res, req)
	}
}

func GetWithAuth(handler http.HandlerFunc) http.HandlerFunc {
	return Get(WithAuth(handler))
}

func PostWithAuth(handler http.HandlerFunc) http.HandlerFunc {
	return Post(WithAuth(handler))
}
