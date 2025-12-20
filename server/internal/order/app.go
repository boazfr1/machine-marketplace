package order

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	db "machine-marketplace/internal/DB/generated"
	"machine-marketplace/internal/machine"
	"machine-marketplace/internal/middleware"
	"machine-marketplace/pkg/kafka"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/labstack/echo/v4"
)

type (
	Module struct {
		P             string
		DB            *db.Queries
		E             *echo.Echo
		KafkaProducer *kafka.Producer
	}

	BuyMachineRequest struct {
		MachineID         int32 `json:"machine_id"`
		DealDurationHours int   `json:"deal_duration_hours"` // Duration in hours
	}
)

var (
	l = slog.New(slog.NewTextHandler(os.Stdout, nil))
)

func (m *Module) SetupOrderRoutes() {
	l.Info("SetupOrderRoutes - setting up order routes")
	g := m.E.Group("/api/v1/order")
	g.Use(middleware.EchoAuth)

	g.GET("", m.ListOfFreeMachines)
	g.POST("/create", m.CreateMachine)
	g.POST("/buy", m.BuyMachine)
	g.GET("/owned-machines", m.GetOwnedMachines)
	g.GET("/bought-machines", m.GetBoughtMachines)
	l.Info("SetupOrderRoutes - order routes setup completed")
}

func (m *Module) ListOfFreeMachines(c echo.Context) error {
	cpuStr := c.QueryParam("cpu")
	ramStr := c.QueryParam("ram")
	gpuStr := c.QueryParam("gpu")

	l.Info("ListOfFreeMachines - fetching machines", "cpu", cpuStr, "ram", ramStr, "gpu", gpuStr)

	machines, err := m.ListOfFreeMachinesService(c.Request().Context(), cpuStr, ramStr, gpuStr)
	if err != nil {
		l.Error("ListOfFreeMachines - failed to get machines", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to get machines")
	}

	l.Info("ListOfFreeMachines - successfully retrieved machines", "count", len(machines))
	return c.JSON(http.StatusOK, machines)
}

func (m *Module) CreateMachine(c echo.Context) error {
	claims := c.Get(string(middleware.ClaimsContextKey)).(*jwt.StandardClaims)

	var params machine.MachineParams
	if err := json.NewDecoder(c.Request().Body).Decode(&params); err != nil {
		l.Error("CreateMachine - invalid request body", "error", err)
		return c.String(http.StatusBadRequest, "Invalid request body")
	}

	if params.Name == "" || params.Ram == 0 || params.Cpu == 0 || params.Memory == 0 || params.Key == "" || params.Host == "" || params.SshUser == "" {
		l.Error("CreateMachine - missing required fields", "name", params.Name, "host", params.Host)
		return c.String(http.StatusBadRequest, "Name, ram, cpu, gpu, memory, key, host, and ssh_user are required")
	}

	num, err := strconv.Atoi(claims.Issuer)
	if err != nil {
		l.Error("CreateMachine - invalid user ID in claims", "error", err)
		return c.String(http.StatusUnauthorized, "Unauthorized")
	}
	ownerID := int32(num)

	l.Info("CreateMachine - creating machine", "name", params.Name, "owner_id", ownerID, "host", params.Host)

	newMachine, err := m.CreateMachineService(c.Request().Context(), params, ownerID)
	if err != nil {
		l.Error("CreateMachine - failed to create machine", "error", err, "name", params.Name)
		return c.String(http.StatusInternalServerError, "Failed to create machine")
	}

	l.Info("CreateMachine - machine created successfully", "machine_id", newMachine.ID, "name", newMachine.Name)
	return c.JSON(http.StatusOK, newMachine)
}

func (m *Module) GetOwnedMachines(c echo.Context) error {
	claims := c.Get(string(middleware.ClaimsContextKey)).(*jwt.StandardClaims)

	ownerID, err := strconv.Atoi(claims.Issuer)
	if err != nil {
		l.Error("GetOwnedMachines - invalid user ID in claims", "error", err)
		return c.String(http.StatusUnauthorized, "Unauthorized")
	}

	l.Info("GetOwnedMachines - fetching owned machines", "owner_id", ownerID)

	machines, err := m.GetOwnedMachinesService(c.Request().Context(), int32(ownerID))
	if err != nil {
		l.Error("GetOwnedMachines - failed to get machines list", "error", err, "owner_id", ownerID)
		return c.String(http.StatusInternalServerError, "Failed to get machines list")
	}

	l.Info("GetOwnedMachines - successfully retrieved owned machines", "owner_id", ownerID, "count", len(machines))
	return c.JSON(http.StatusOK, machines)
}

func (m *Module) GetBoughtMachines(c echo.Context) error {
	claims := c.Get(string(middleware.ClaimsContextKey)).(*jwt.StandardClaims)

	num, err := strconv.Atoi(claims.Issuer)
	if err != nil {
		l.Error("GetBoughtMachines - invalid user ID in claims", "error", err)
		return c.String(http.StatusUnauthorized, "Unauthorized")
	}

	l.Info("GetBoughtMachines - fetching bought machines", "buyer_id", num)

	machines, err := m.GetBoughtMachinesService(c.Request().Context(), int32(num))
	if err != nil {
		l.Error("GetBoughtMachines - failed to get machines list", "error", err, "buyer_id", num)
		return c.String(http.StatusInternalServerError, "Failed to get machines list")
	}

	l.Info("GetBoughtMachines - successfully retrieved bought machines", "buyer_id", num, "count", len(machines))
	return c.JSON(http.StatusOK, machines)
}

func (m *Module) BuyMachine(c echo.Context) error {
	claims := c.Get(string(middleware.ClaimsContextKey)).(*jwt.StandardClaims)

	var req BuyMachineRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		l.Error("BuyMachine - invalid request body", "error", err)
		return c.String(http.StatusBadRequest, "Invalid request body")
	}

	if req.MachineID == 0 {
		l.Error("BuyMachine - missing machine_id")
		return c.String(http.StatusBadRequest, "machine_id is required")
	}

	if req.DealDurationHours <= 0 {
		req.DealDurationHours = 24
		l.Info("BuyMachine - using default deal duration", "duration_hours", req.DealDurationHours)
	}

	buyerID, err := strconv.Atoi(claims.Issuer)
	if err != nil {
		l.Error("BuyMachine - invalid user ID in claims", "error", err)
		return c.String(http.StatusUnauthorized, "Unauthorized")
	}

	l.Info("BuyMachine - processing purchase", "machine_id", req.MachineID, "buyer_id", buyerID, "duration_hours", req.DealDurationHours)

	result, err := m.BuyMachineService(c.Request().Context(), req.MachineID, int32(buyerID), req.DealDurationHours)
	if err != nil {
		if err == sql.ErrNoRows {
			l.Error("BuyMachine - machine already purchased", "machine_id", req.MachineID, "buyer_id", buyerID)
			return c.String(http.StatusConflict, "Machine is already purchased")
		}
		l.Error("BuyMachine - failed to purchase machine", "error", err, "machine_id", req.MachineID, "buyer_id", buyerID)
		return c.String(http.StatusInternalServerError, "Failed to purchase machine")
	}

	l.Info("BuyMachine - machine purchased successfully", "machine_id", req.MachineID, "buyer_id", buyerID, "expiration", result.DealExpiration)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":         "Machine purchased successfully",
		"machine":         result.Machine,
		"deal_expiration": result.DealExpiration,
	})
}
