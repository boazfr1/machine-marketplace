package user

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	db "machine-marketplace/internal/DB/generated"
	"machine-marketplace/pkg/database"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func Login(res http.ResponseWriter, req *http.Request) {
	var params struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(req.Body).Decode(&params); err != nil {
		fmt.Println("Invalid request body")
		http.Error(res, "Invalid request body", http.StatusBadRequest)
		return
	}

	if params.Email == "" || params.Password == "" {
		fmt.Println("Email and password are required")
		http.Error(res, "Email and password are required", http.StatusBadRequest)
		return
	}

	userByEmail, err := database.Queries.GetUserByEmail(req.Context(), params.Email)
	if err != nil {
		fmt.Println("User not found:", err)
		res.Header().Set("Content-Type", "application/json")
		json.NewEncoder(res).Encode(map[string]string{
			"message": "user not found",
		})
		return
	}

	// Remove the \x prefix and decode the hex string
	passwordString := string(userByEmail.Password)
	if strings.HasPrefix(passwordString, `\x`) {
		passwordString = passwordString[2:] // Remove \x prefix
		decodedBytes, err := hex.DecodeString(passwordString)
		if err != nil {
			fmt.Println("Failed to decode password:", err)
			res.Header().Set("Content-Type", "application/json")
			json.NewEncoder(res).Encode(map[string]string{
				"message": "Internal server error",
			})
			return
		}
		userByEmail.Password = decodedBytes
	}

	if err := bcrypt.CompareHashAndPassword(userByEmail.Password, []byte(params.Password)); err != nil {
		fmt.Println("Password comparison failed:", err)
		res.Header().Set("Content-Type", "application/json")
		json.NewEncoder(res).Encode(map[string]string{
			"message": "user or password are incorrect",
		})
		return
	}

	response := map[string]interface{}{
		"message": "Login successful",
		"user": map[string]string{
			"name":  userByEmail.Name,
			"email": userByEmail.Email,
		},
	}

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(response)
}

func SignUp(res http.ResponseWriter, req *http.Request) {
	var params struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

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
