package user

import (
	"net/http"
	"time"
)

func createAuthCookie(token string) *http.Cookie {
	return &http.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	}
}
