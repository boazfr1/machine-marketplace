package user

import (
	"encoding/json"
	db "machine-marketplace/internal/DB/generated"
	"machine-marketplace/pkg/database"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func SignUp(res http.ResponseWriter, req *http.Request) {
	var params userParams
	
	if err := json.NewDecoder(req.Body).Decode(&params); err != nil {
		http.Error(res, "Invalid request body", http.StatusBadRequest)
		return
	}

	if params.Name == "" || params.Email == "" || params.Password == "" {
		http.Error(res, "Name, email, and password are required", http.StatusBadRequest)
		return
	}

	cryptPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), 14)
	if err != nil {
		http.Error(res, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	createParams := db.CreateUserParams{
		Name:    params.Name,
		Email:   params.Email,
		Column3: cryptPassword,
	}

	user, err := database.Queries.CreateUser(req.Context(), createParams)
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			http.Error(res, "Email already exists", http.StatusConflict)
			return
		}
		http.Error(res, "Failed to create user", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "User created successfully",
		"user": map[string]interface{}{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	}

	res.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(response); err != nil {
		http.Error(res, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
