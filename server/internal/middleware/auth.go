package middleware

import (
	"context"
	user "machine-marketplace/internal/user"
	"net/http"
)

type contextKey string

const ClaimsContextKey contextKey = "claims"

func WithAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie("jwt")
		if err != nil {
			http.Error(res, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, err := user.ValidateToken(cookie.Value)
		if err != nil {
			http.Error(res, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(req.Context(), ClaimsContextKey, claims)
		next.ServeHTTP(res, req.WithContext(ctx))
	}
}
