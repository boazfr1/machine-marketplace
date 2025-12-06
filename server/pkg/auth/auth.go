package auth

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"golang.org/x/crypto/bcrypt"
)

const secretKey = "secret"

func GenerateToken(userID int32) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(userID)),
		ExpiresAt: jwt.At(time.Now().Add(time.Hour * 24)),
	})

	return claims.SignedString([]byte(secretKey))
}

func ValidateToken(tokenString string) (*jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	return token.Claims.(*jwt.StandardClaims), nil
}

func VerifyPassword(storedPassword []byte, passwordString string) error {
	if strings.HasPrefix(passwordString, `\x`) {
		passwordString = passwordString[2:]
		decodedBytes, err := hex.DecodeString(passwordString)
		if err != nil {
			return fmt.Errorf("failed to decode password: %w", err)
		}
		storedPassword = decodedBytes
	}
	return bcrypt.CompareHashAndPassword(storedPassword, []byte(passwordString))
}

func CreateAuthCookie(token string) *http.Cookie {
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


