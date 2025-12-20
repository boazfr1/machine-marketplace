package middleware

import (
	"log/slog"
	"net/http"
	"os"
)

var (
	corsLogger = slog.New(slog.NewTextHandler(os.Stdout, nil))
)

func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			corsLogger.Info("EnableCORS - setting origin", "origin", origin, "path", r.URL.Path)
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			corsLogger.Info("EnableCORS - handling OPTIONS preflight", "origin", origin, "path", r.URL.Path)
			w.WriteHeader(http.StatusOK)
			return
		}

		corsLogger.Info("EnableCORS - processing request", "method", r.Method, "path", r.URL.Path, "origin", origin)
		next.ServeHTTP(w, r)
	})
}
