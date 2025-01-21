package user

import (
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
)

const secretKey = "secret"

func generateToken(userID int32) (string, error) {
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
