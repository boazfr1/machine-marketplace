package user

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	db "machine-marketplace/internal/DB/generated"
	data "machine-marketplace/internal/data"
	"machine-marketplace/pkg/auth"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type (
	Module struct {
		P string
		DB *db.Queries
		E *echo.Echo
	}
)

func (m *Module) SetupUserRoutes() {
	m.E.GET("/api/v1/health", m.health)
	m.E.POST("/api/v1/sign-up", m.SignUp)
	m.E.POST("/api/v1/login", m.Login)
	m.E.POST("/api/v1/logout", m.Logout)
	m.E.GET("/api/v1/user", m.User)
}

func (m *Module) Login(c echo.Context) error {
	var params data.UserParams
	if err := json.NewDecoder(c.Request().Body).Decode(&params); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request body")
	}

	if params.Email == "" || params.Password == "" {
		return c.String(http.StatusBadRequest, "Email and password are required")
	}

	userByEmail, err := m.DB.GetUserByEmail(c.Request().Context(), params.Email)
	if err != nil {
		return c.String(http.StatusUnauthorized, "User not found")
	}

	if err := auth.VerifyPassword(userByEmail.Password, params.Password); err != nil {
		return c.String(http.StatusUnauthorized, "Invalid credentials")
	}

	token, err := auth.GenerateToken(userByEmail.ID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to generate token")
	}

	c.SetCookie(auth.CreateAuthCookie(token))

	c.Response().Header().Set("Content-Type", "application/json")
	return json.NewEncoder(c.Response()).Encode(map[string]interface{}{
		"message": "Login successful",
	})
}

func (m *Module) User(c echo.Context) error {
	cookie, err := c.Cookie("jwt")
	if err != nil {
		return c.String(http.StatusUnauthorized, "Unauthorized")
	}

	claims, err := auth.ValidateToken(cookie.Value)
	if err != nil {
		return c.String(http.StatusUnauthorized, "Unauthorized")
	}

	c.Response().Header().Set("Content-Type", "application/json")
	return json.NewEncoder(c.Response()).Encode(map[string]interface{}{
		"issuer": claims.Issuer,
	})
}

func (m *Module) health(c echo.Context) error {
	data := map[string]string{
		"massage": "be kind",
	}

	c.Response().Header().Set("Content-Type", "application/json")
	return json.NewEncoder(c.Response()).Encode(data)
}

func (m *Module) Logout(c echo.Context) error {
	cookie := auth.CreateAuthCookie("")
	cookie.Expires = time.Now().Add(-time.Hour)
	c.SetCookie(cookie)

	c.Response().Header().Set("Content-Type", "application/json")
	return json.NewEncoder(c.Response()).Encode(map[string]interface{}{
		"message": "Logout successful",
	})
}

func (m *Module) SignUp(c echo.Context) error {
	var params data.UserParams

	if err := json.NewDecoder(c.Request().Body).Decode(&params); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request body")
	}

	if params.Name == "" || params.Email == "" || params.Password == "" {
		return c.String(http.StatusBadRequest, "Name, email, and password are required")
	}

	cryptPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), 14)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to hash password")
	}

	createParams := db.CreateUserParams{
		Name:    params.Name,
		Email:   params.Email,
		Column3: cryptPassword,
	}

	user, err := m.DB.CreateUser(c.Request().Context(), createParams)
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return c.String(http.StatusConflict, "Email already exists")
		}
		return c.String(http.StatusInternalServerError, "Failed to create user")
	}

	response := map[string]interface{}{
		"message": "User created successfully",
		"user": map[string]interface{}{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	}

	c.Response().Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(c.Response()).Encode(response); err != nil {
		return c.String(http.StatusInternalServerError, "Failed to encode response")
	}
	return nil
}
