package user

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
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
		P  string
		DB *db.Queries
		E  *echo.Echo
	}
)

var (
	l = slog.New(slog.NewTextHandler(os.Stdout, nil))
)

func (m *Module) SetupUserRoutes() {
	l.Info("SetupUserRoutes - setting up user routes")
	m.E.GET("/api/v1/health", m.health)
	m.E.POST("/api/v1/sign-up", m.SignUp)
	m.E.POST("/api/v1/login", m.Login)
	m.E.POST("/api/v1/logout", m.Logout)
	m.E.GET("/api/v1/user", m.User)
	l.Info("SetupUserRoutes - user routes setup completed")
}

func (m *Module) Login(c echo.Context) error {
	var params data.UserParams
	if err := json.NewDecoder(c.Request().Body).Decode(&params); err != nil {
		l.Error("Login - invalid request body", "error", err)
		return c.String(http.StatusBadRequest, "Invalid request body")
	}

	if params.Email == "" || params.Password == "" {
		l.Error("Login - missing email or password", "email", params.Email)
		return c.String(http.StatusBadRequest, "Email and password are required")
	}

	l.Info("Login - attempting login", "email", params.Email)

	userByEmail, err := m.DB.GetUserByEmail(c.Request().Context(), params.Email)
	if err != nil {
		l.Error("Login - user not found", "error", err, "email", params.Email)
		return c.String(http.StatusUnauthorized, "User not found")
	}

	if err := auth.VerifyPassword(userByEmail.Password, params.Password); err != nil {
		l.Error("Login - invalid credentials", "error", err, "email", params.Email)
		return c.String(http.StatusUnauthorized, "Invalid credentials")
	}

	token, err := auth.GenerateToken(userByEmail.ID)
	if err != nil {
		l.Error("Login - failed to generate token", "error", err, "user_id", userByEmail.ID)
		return c.String(http.StatusInternalServerError, "Failed to generate token")
	}

	c.SetCookie(auth.CreateAuthCookie(token))

	l.Info("Login - login successful", "user_id", userByEmail.ID, "email", params.Email)
	c.Response().Header().Set("Content-Type", "application/json")
	return json.NewEncoder(c.Response()).Encode(map[string]interface{}{
		"message": "Login successful",
	})
}

func (m *Module) User(c echo.Context) error {
	cookie, err := c.Cookie("jwt")
	if err != nil {
		l.Error("User - no JWT cookie found", "error", err)
		return c.String(http.StatusUnauthorized, "Unauthorized")
	}

	claims, err := auth.ValidateToken(cookie.Value)
	if err != nil {
		l.Error("User - invalid token", "error", err)
		return c.String(http.StatusUnauthorized, "Unauthorized")
	}

	l.Info("User - user info retrieved", "user_id", claims.Issuer)
	c.Response().Header().Set("Content-Type", "application/json")
	return json.NewEncoder(c.Response()).Encode(map[string]interface{}{
		"issuer": claims.Issuer,
	})
}

func (m *Module) health(c echo.Context) error {
	l.Info("health - health check requested")
	data := map[string]string{
		"massage": "be kind",
	}

	c.Response().Header().Set("Content-Type", "application/json")
	return json.NewEncoder(c.Response()).Encode(data)
}

func (m *Module) Logout(c echo.Context) error {
	l.Info("Logout - user logout requested")
	cookie := auth.CreateAuthCookie("")
	cookie.Expires = time.Now().Add(-time.Hour)
	c.SetCookie(cookie)

	l.Info("Logout - logout successful")
	c.Response().Header().Set("Content-Type", "application/json")
	return json.NewEncoder(c.Response()).Encode(map[string]interface{}{
		"message": "Logout successful",
	})
}

func (m *Module) SignUp(c echo.Context) error {
	var params data.UserParams

	if err := json.NewDecoder(c.Request().Body).Decode(&params); err != nil {
		l.Error("SignUp - invalid request body", "error", err)
		return c.String(http.StatusBadRequest, "Invalid request body")
	}

	if params.Name == "" || params.Email == "" || params.Password == "" {
		l.Error("SignUp - missing required fields", "name", params.Name, "email", params.Email)
		return c.String(http.StatusBadRequest, "Name, email, and password are required")
	}

	l.Info("SignUp - creating new user", "name", params.Name, "email", params.Email)

	cryptPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), 14)
	if err != nil {
		l.Error("SignUp - failed to hash password", "error", err, "email", params.Email)
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
			l.Error("SignUp - email already exists", "error", err, "email", params.Email)
			return c.String(http.StatusConflict, "Email already exists")
		}
		l.Error("SignUp - failed to create user", "error", err, "email", params.Email)
		return c.String(http.StatusInternalServerError, "Failed to create user")
	}

	l.Info("SignUp - user created successfully", "user_id", user.ID, "name", user.Name, "email", user.Email)

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
		l.Error("SignUp - failed to encode response", "error", err, "user_id", user.ID)
		return c.String(http.StatusInternalServerError, "Failed to encode response")
	}
	return nil
}
