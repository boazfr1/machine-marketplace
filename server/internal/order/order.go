package order

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	db "machine-marketplace/internal/DB/generated"
	"machine-marketplace/internal/machine"
	"machine-marketplace/internal/middleware"
	"machine-marketplace/pkg/kafka"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/labstack/echo/v4"
)

type (
	Module struct {
		P              string
		DB             *db.Queries
		E              *echo.Echo
		KafkaProducer  *kafka.Producer
	}
)

func (m *Module) SetupOrderRoutes() {
	g := m.E.Group("/api/v1/order")
	g.Use(middleware.EchoAuth)

	g.GET("", m.ListOfFreeMachines)
	g.POST("/create", m.CreateMachine)
	g.POST("/buy", m.BuyMachine)
	g.GET("/owned-machines", m.GetOwnedMachines)
	g.GET("/bought-machines", m.GetBoughtMachines)
}

func (m *Module) ListOfFreeMachines(c echo.Context) error {
	// Check for filter parameters
	cpuStr := c.QueryParam("cpu")
	ramStr := c.QueryParam("ram")
	gpuStr := c.QueryParam("gpu")

	// If any filter is present, use the filter query
	if cpuStr != "" || ramStr != "" || gpuStr != "" {
		var cpu, ram, gpu sql.NullInt32

		if cpuStr != "" {
			cpuVal, err := strconv.Atoi(cpuStr)
			if err == nil {
				cpu = sql.NullInt32{Int32: int32(cpuVal), Valid: true}
			}
		}

		if ramStr != "" {
			ramVal, err := strconv.Atoi(ramStr)
			if err == nil {
				ram = sql.NullInt32{Int32: int32(ramVal), Valid: true}
			}
		}

		if gpuStr != "" {
			gpuVal, err := strconv.Atoi(gpuStr)
			if err == nil {
				gpu = sql.NullInt32{Int32: int32(gpuVal), Valid: true}
			}
		}

		machines, err := m.DB.FilterAvailableMachines(c.Request().Context(), db.FilterAvailableMachinesParams{
			Column1: cpu,
			Column2: ram,
			Column3: gpu,
		})
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to filter machines")
		}
		return c.JSON(http.StatusOK, machines)
	}

	// No filters, return all available machines
	machines, err := m.DB.ListAvailableMachines(c.Request().Context())
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to get machines")
	}
	return c.JSON(http.StatusOK, machines)
}

func (m *Module) CreateMachine(c echo.Context) error {
	claims := c.Get(string(middleware.ClaimsContextKey)).(*jwt.StandardClaims)

	var params machine.MachineParams
	if err := json.NewDecoder(c.Request().Body).Decode(&params); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request body")
	}

	if params.Name == "" || params.Ram == 0 || params.Cpu == 0 || params.Memory == 0 || params.Key == "" || params.Host == "" || params.SshUser == "" {
		return c.String(http.StatusBadRequest, "Name, ram, cpu, gpu, memory, key, host, and ssh_user are required")
	}

	num, err := strconv.Atoi(claims.Issuer)
	if err != nil {
		return c.String(http.StatusUnauthorized, "Unauthorized")
	}
	ownerID := int32(num)

	err = machine.TriedToConnectForFirstTime(params.Host, params.SshUser, params.Key)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to connect to machine")
	}

	createParams := db.CreateMachineParams{
		Name:    params.Name,
		Ram:     params.Ram,
		Cpu:     params.Cpu,
		Gpu:     params.Gpu,
		Memory:  params.Memory,
		Key:     sql.NullString{String: params.Key, Valid: true},
		OwnerID: ownerID,
		Host:    params.Host,
		SshUser: params.SshUser,
	}

	newMachine, err := m.DB.CreateMachine(c.Request().Context(), createParams)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create machine")
	}

	return c.JSON(http.StatusOK, newMachine)
}

func (m *Module) GetOwnedMachines(c echo.Context) error {
	claims := c.Get(string(middleware.ClaimsContextKey)).(*jwt.StandardClaims)

	ownerID, err := strconv.Atoi(claims.Issuer)
	if err != nil {
		return c.String(http.StatusUnauthorized, "Unauthorized")
	}

	machines, err := m.DB.ListMachinesByOwnerID(c.Request().Context(), int32(ownerID))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to get machines list")
	}
	return c.JSON(http.StatusOK, machines)
}

func (m *Module) GetBoughtMachines(c echo.Context) error {
	claims := c.Get(string(middleware.ClaimsContextKey)).(*jwt.StandardClaims)

	num, err := strconv.Atoi(claims.Issuer)
	if err != nil {
		return c.String(http.StatusUnauthorized, "Unauthorized")
	}

	buyerID := sql.NullInt32{
		Int32: int32(num),
		Valid: true,
	}

	machines, err := m.DB.ListMachinesByBuyerID(c.Request().Context(), buyerID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to get machines list")
	}
	return c.JSON(http.StatusOK, machines)
}

type BuyMachineRequest struct {
	MachineID      int32 `json:"machine_id"`
	DealDurationHours int `json:"deal_duration_hours"` // Duration in hours
}

func (m *Module) BuyMachine(c echo.Context) error {
	claims := c.Get(string(middleware.ClaimsContextKey)).(*jwt.StandardClaims)

	var req BuyMachineRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request body")
	}

	if req.MachineID == 0 {
		return c.String(http.StatusBadRequest, "machine_id is required")
	}

	if req.DealDurationHours <= 0 {
		req.DealDurationHours = 24 // Default to 24 hours
	}

	buyerID, err := strconv.Atoi(claims.Issuer)
	if err != nil {
		return c.String(http.StatusUnauthorized, "Unauthorized")
	}

	// Get machine details to verify it exists and is available
	machineDetails, err := m.DB.GetMachineByID(c.Request().Context(), req.MachineID)
	if err != nil {
		return c.String(http.StatusNotFound, "Machine not found")
	}

	// Check if machine is already purchased
	if machineDetails.BuyerID.Valid {
		return c.String(http.StatusConflict, "Machine is already purchased")
	}

	// Update machine with buyer_id
	updateParams := db.UpdateMachineBuyerParams{
		BuyerID: sql.NullInt32{
			Int32: int32(buyerID),
			Valid: true,
		},
		Key: machineDetails.Key,
		ID:  req.MachineID,
	}

	updatedMachine, err := m.DB.UpdateMachineBuyer(c.Request().Context(), updateParams)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to purchase machine")
	}

	// Send purchase event to Kafka
	dealExpiration := time.Now().Add(time.Duration(req.DealDurationHours) * time.Hour)
	purchaseEvent := kafka.PurchaseEvent{
		MachineID:      req.MachineID,
		BuyerID:        int32(buyerID),
		DealExpiration: dealExpiration,
		PurchaseTime:   time.Now(),
		MachineName:    updatedMachine.Name,
	}

	if err := m.KafkaProducer.SendPurchaseEvent(c.Request().Context(), purchaseEvent); err != nil {
		// Log error but don't fail the purchase
		fmt.Printf("Failed to send purchase event to Kafka: %v\n", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":         "Machine purchased successfully",
		"machine":         updatedMachine,
		"deal_expiration": dealExpiration,
	})
}
