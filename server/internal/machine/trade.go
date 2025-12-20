package machine

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	db "machine-marketplace/internal/DB/generated"
	middleware "machine-marketplace/internal/middleware"
	database "machine-marketplace/pkg/database"
	"net/http"
	"os"
	"strconv"

	"github.com/dgrijalva/jwt-go/v4"
)

var (
	tradeLogger = slog.New(slog.NewTextHandler(os.Stdout, nil))
)

type MachineParams struct {
	Name    string `json:"name"`
	Ram     int32  `json:"ram"`
	Cpu     int32  `json:"cpu"`
	Gpu     int32  `json:"gpu"`
	Memory  int32  `json:"memory"`
	Key     string `json:"key"`
	Host    string `json:"host"`
	SshUser string `json:"ssh_user"`
}

func CreateMachine(res http.ResponseWriter, req *http.Request) {
	claims := req.Context().Value(middleware.ClaimsContextKey).(*jwt.StandardClaims)

	var params MachineParams
	if err := json.NewDecoder(req.Body).Decode(&params); err != nil {
		tradeLogger.Error("CreateMachine - invalid request body", "error", err)
		http.Error(res, "Invalid request body", http.StatusBadRequest)
		return
	}

	if params.Name == "" || params.Ram == 0 || params.Cpu == 0 || params.Memory == 0 || params.Key == "" || params.Host == "" || params.SshUser == "" {
		tradeLogger.Error("CreateMachine - missing required fields", "name", params.Name, "host", params.Host)
		http.Error(res, "Name, ram, cpu, gpu, memory, key, host, and ssh_user are required", http.StatusBadRequest)
		return
	}

	num, err := strconv.Atoi(claims.Issuer)
	if err != nil {
		tradeLogger.Error("CreateMachine - invalid user ID in claims", "error", err)
		http.Error(res, "Unauthorized", http.StatusUnauthorized)
		return
	}
	ownerID := int32(num)

	tradeLogger.Info("CreateMachine - creating machine", "name", params.Name, "owner_id", ownerID, "host", params.Host)

	err = TriedToConnectForFirstTime(params.Host, params.SshUser, params.Key)
	if err != nil {
		tradeLogger.Error("CreateMachine - failed to connect to machine", "error", err, "host", params.Host)
		http.Error(res, "Failed to connect to machine", http.StatusInternalServerError)
		return
	}

	createParams := db.CreateMachineParams{
		Name:    params.Name,
		Ram:     params.Ram,
		Cpu:     params.Cpu,
		Gpu:     params.Gpu,
		Memory:  params.Memory,
		Key:     sql.NullString{String: params.Key, Valid: true},
		OwnerID: ownerID,
	}

	machine, err := database.Queries.CreateMachine(req.Context(), createParams)
	if err != nil {
		tradeLogger.Error("CreateMachine - failed to create machine", "error", err, "name", params.Name)
		http.Error(res, "Failed to create machine", http.StatusInternalServerError)
		return
	}

	tradeLogger.Info("CreateMachine - machine created successfully", "machine_id", machine.ID, "name", machine.Name)
	json.NewEncoder(res).Encode(machine)
}

func ListOfFreeMachines(res http.ResponseWriter, req *http.Request) {
	tradeLogger.Info("ListOfFreeMachines - fetching available machines")
	
	machines, err := database.Queries.ListAvailableMachines(req.Context())
	if err != nil {
		tradeLogger.Error("ListOfFreeMachines - failed to get machines", "error", err)
		http.Error(res, "Failed to get machines", http.StatusInternalServerError)
		return
	}
	
	tradeLogger.Info("ListOfFreeMachines - machines retrieved successfully", "count", len(machines))
	json.NewEncoder(res).Encode(machines)
}

func GetMachineByID(res http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	if id == "" {
		tradeLogger.Error("GetMachineByID - missing ID parameter")
		http.Error(res, "ID is required", http.StatusBadRequest)
		return
	}

	machineID, err := strconv.Atoi(id)
	if err != nil {
		tradeLogger.Error("GetMachineByID - invalid ID", "error", err, "id", id)
		http.Error(res, "Invalid ID", http.StatusBadRequest)
		return
	}

	tradeLogger.Info("GetMachineByID - fetching machine", "machine_id", machineID)

	machine, err := database.Queries.GetMachineByID(req.Context(), int32(machineID))
	if err != nil {
		tradeLogger.Error("GetMachineByID - failed to get machine", "error", err, "machine_id", machineID)
		http.Error(res, "Failed to get machine", http.StatusInternalServerError)
		return
	}
	
	tradeLogger.Info("GetMachineByID - machine retrieved successfully", "machine_id", machineID, "name", machine.Name)
	json.NewEncoder(res).Encode(machine)
}

func GetOwnedMachines(res http.ResponseWriter, req *http.Request) {
	claims := req.Context().Value(middleware.ClaimsContextKey).(*jwt.StandardClaims)

	ownerID, err := strconv.Atoi(claims.Issuer)
	if err != nil {
		tradeLogger.Error("GetOwnedMachines - invalid user ID in claims", "error", err)
		http.Error(res, "Unauthorized", http.StatusUnauthorized)
		return
	}

	tradeLogger.Info("GetOwnedMachines - fetching owned machines", "owner_id", ownerID)

	machines, err := database.Queries.ListMachinesByOwnerID(req.Context(), int32(ownerID))
	if err != nil {
		tradeLogger.Error("GetOwnedMachines - failed to get machines list", "error", err, "owner_id", ownerID)
		http.Error(res, "Failed to get machines list", http.StatusInternalServerError)
		return
	}
	
	tradeLogger.Info("GetOwnedMachines - owned machines retrieved", "owner_id", ownerID, "count", len(machines))
	json.NewEncoder(res).Encode(machines)
}

func GetBoughtMachines(res http.ResponseWriter, req *http.Request) {
	claims := req.Context().Value(middleware.ClaimsContextKey).(*jwt.StandardClaims)

	num, err := strconv.Atoi(claims.Issuer)
	if err != nil {
		tradeLogger.Error("GetBoughtMachines - invalid user ID in claims", "error", err)
		http.Error(res, "Unauthorized", http.StatusUnauthorized)
		return
	}

	tradeLogger.Info("GetBoughtMachines - fetching bought machines", "buyer_id", num)

	buyerID := sql.NullInt32{
		Int32: int32(num),
		Valid: true,
	}

	machines, err := database.Queries.ListMachinesByBuyerID(req.Context(), buyerID)
	if err != nil {
		tradeLogger.Error("GetBoughtMachines - failed to get machines list", "error", err, "buyer_id", num)
		http.Error(res, "Failed to get machines list", http.StatusInternalServerError)
		return
	}
	
	tradeLogger.Info("GetBoughtMachines - bought machines retrieved", "buyer_id", num, "count", len(machines))
	json.NewEncoder(res).Encode(machines)
}

func BuyMachine(res http.ResponseWriter, req *http.Request) {

}
