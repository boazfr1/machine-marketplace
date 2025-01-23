package user

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"machine-marketplace/pkg/database"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type userParams struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func verifyPassword(storedPassword []byte, passwordString string) error {
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

func Login(res http.ResponseWriter, req *http.Request) {
	var params userParams
	if err := json.NewDecoder(req.Body).Decode(&params); err != nil {
		http.Error(res, "Invalid request body", http.StatusBadRequest)
		return
	}

	if params.Email == "" || params.Password == "" {
		http.Error(res, "Email and password are required", http.StatusBadRequest)
		return
	}

	userByEmail, err := database.Queries.GetUserByEmail(req.Context(), params.Email)
	if err != nil {
		http.Error(res, "User not found", http.StatusUnauthorized)
		return
	}

	if err := verifyPassword(userByEmail.Password, params.Password); err != nil {
		http.Error(res, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := generateToken(userByEmail.ID)
	if err != nil {
		http.Error(res, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(res, createAuthCookie(token))

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(map[string]interface{}{
		"message": "Login successful",
	})
}

func User(res http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("jwt")
	if err != nil {
		http.Error(res, "Unauthorized", http.StatusUnauthorized)
		return
	}

	claims, err := ValidateToken(cookie.Value)
	if err != nil {
		http.Error(res, "Unauthorized", http.StatusUnauthorized)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(map[string]interface{}{
		"issuer": claims.Issuer,
	})
}
