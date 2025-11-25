package middleware

import (
	"context"
	"machine-marketplace/pkg/auth"
	"net/http"

	"github.com/labstack/echo/v4"
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

		claims, err := auth.ValidateToken(cookie.Value)
		if err != nil {
			http.Error(res, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(req.Context(), ClaimsContextKey, claims)
		next.ServeHTTP(res, req.WithContext(ctx))
	}
}

func EchoAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("jwt")
		if err != nil {
			return c.String(http.StatusUnauthorized, "Unauthorized")
		}

		claims, err := auth.ValidateToken(cookie.Value)
		if err != nil {
			return c.String(http.StatusUnauthorized, "Unauthorized")
		}

		c.Set(string(ClaimsContextKey), claims)
		return next(c)
	}
}
