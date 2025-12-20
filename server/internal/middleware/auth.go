package middleware

import (
	"context"
	"log/slog"
	"machine-marketplace/pkg/auth"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

var (
	l = slog.New(slog.NewTextHandler(os.Stdout, nil))
)

type contextKey string

const ClaimsContextKey contextKey = "claims"

func WithAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie("jwt")
		if err != nil {
			l.Error("WithAuth - no JWT cookie found", "error", err, "path", req.URL.Path)
			http.Error(res, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, err := auth.ValidateToken(cookie.Value)
		if err != nil {
			l.Error("WithAuth - invalid token", "error", err, "path", req.URL.Path)
			http.Error(res, "Unauthorized", http.StatusUnauthorized)
			return
		}

		l.Info("WithAuth - authentication successful", "user_id", claims.Issuer, "path", req.URL.Path)
		ctx := context.WithValue(req.Context(), ClaimsContextKey, claims)
		next.ServeHTTP(res, req.WithContext(ctx))
	}
}

func EchoAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("jwt")
		if err != nil {
			l.Error("EchoAuth - no JWT cookie found", "error", err, "path", c.Request().URL.Path)
			return c.String(http.StatusUnauthorized, "Unauthorized")
		}

		claims, err := auth.ValidateToken(cookie.Value)
		if err != nil {
			l.Error("EchoAuth - invalid token", "error", err, "path", c.Request().URL.Path)
			return c.String(http.StatusUnauthorized, "Unauthorized")
		}

		l.Info("EchoAuth - authentication successful", "user_id", claims.Issuer, "path", c.Request().URL.Path)
		c.Set(string(ClaimsContextKey), claims)
		return next(c)
	}
}
