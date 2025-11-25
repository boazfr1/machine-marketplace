package order

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	db "machine-marketplace/internal/DB/generated"
	"machine-marketplace/internal/machine"
	"machine-marketplace/internal/middleware"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/labstack/echo/v4"
)

type (
	Module struct {
		P  string
		DB *db.Queries
		E  *echo.Echo
	}
)

func (m *Module) SetupOrderRoutes() {
	g := m.E.Group("/api/v1/order")
	g.Use(middleware.EchoAuth)

	g.GET("", m.ListOfFreeMachines)
	g.POST("/create", m.CreateMachine)
	g.GET("/connect", m.WebSocketHandler)
	g.GET("/owned-machines", m.GetOwnedMachines)
	g.GET("/bought-machines", m.GetBoughtMachines)
}

func (m *Module) ListOfFreeMachines(c echo.Context) error {
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
		return c.String(http.StatusBadRequest, "Name, ram, cpu, memory, key, host, and ssh_user are required")
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
		Memory:  params.Memory,
		Key:     sql.NullString{String: params.Key, Valid: true},
		OwnerID: ownerID,
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

func (m *Module) WebSocketHandler(c echo.Context) error {
	claims := c.Get(string(middleware.ClaimsContextKey)).(*jwt.StandardClaims)

	query := c.Request().URL.Query()

	ownerID, err := strconv.Atoi(query.Get("owner_name"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid owner ID")
	}

	createParams := db.GetMachineByNameAndOwnerParams{
		Name:    query.Get("machine_name"),
		OwnerID: int32(ownerID),
	}

	params, err := m.DB.GetMachineByNameAndOwner(c.Request().Context(), createParams)
	if err != nil {
		fmt.Println("err = ", err)
		return c.String(http.StatusNotFound, "Machine not found")
	}

	num, err := strconv.Atoi(claims.Issuer)
	if err != nil {
		fmt.Println("err = ", err)
		return c.String(http.StatusUnauthorized, "Unauthorized")
	}
	buyerID := int32(num)

	if params.BuyerID.Int32 != buyerID {
		return c.String(http.StatusForbidden, "You are not the owner of this machine")
	}

	wsConn, err := machine.Upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		fmt.Println("Error upgrading:", err)
		return err
	}

	sshClient, err := machine.CreateSSHClient(params.Host, params.SshUser, params.Key.String)
	if err != nil {
		wsConn.Close()
		fmt.Printf("Error creating SSH client: %v\n", err)
		return nil // Connection upgraded, so we return nil or error
	}

	conn := machine.NewConnection(wsConn, sshClient)

	machine.Manager.AddConnection(params.Host, conn)

	go machine.HandleConnection(params.Host, conn)

	return nil
}
