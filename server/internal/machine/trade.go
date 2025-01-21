package machine

import (
	"database/sql"
	"encoding/json"
	db "machine-marketplace/internal/DB/generated"
	user "machine-marketplace/internal/user"
	database "machine-marketplace/pkg/database"
	"strconv"

	"net/http"
)

func CreateMachine(res http.ResponseWriter, req *http.Request) {
	var MachineParams struct {
		Name    string `json:"name"`
		Ram     int32  `json:"ram"`
		Cpu     int32  `json:"cpu"`
		Memory  int32  `json:"memory"`
		Key     string `json:"key"`
		Host    string `json:"host"`
		SshUser string `json:"ssh_user"`
	}

	if err := json.NewDecoder(req.Body).Decode(&MachineParams); err != nil {
		http.Error(res, "Invalid request body", http.StatusBadRequest)
		return
	}

	if MachineParams.Name == "" || MachineParams.Ram == 0 || MachineParams.Cpu == 0 || MachineParams.Memory == 0 || MachineParams.Key == "" || MachineParams.Host == "" || MachineParams.SshUser == "" {
		http.Error(res, "Name, ram, cpu, memory, key, host, and ssh_user are required", http.StatusBadRequest)
		return
	}

	cookie, err := req.Cookie("jwt")
	if err != nil {
		http.Error(res, "Unauthorized", http.StatusUnauthorized)
		return
	}

	claims, err := user.ValidateToken(cookie.Value)
	if err != nil {
		http.Error(res, "Unauthorized", http.StatusUnauthorized)
		return
	}

	num, err := strconv.Atoi(claims.Issuer)
	if err != nil {
		http.Error(res, "Unauthorized", http.StatusUnauthorized)
		return
	}
	ownerID := int32(num)

	err = triedToConnectForFirstTime(MachineParams.Host, MachineParams.SshUser, MachineParams.Key)
	if err != nil {
		http.Error(res, "Failed to connect to machine", http.StatusInternalServerError)
		return
	}

	createParams := db.CreateMachineParams{
		Name:    MachineParams.Name,
		Ram:     MachineParams.Ram,
		Cpu:     MachineParams.Cpu,
		Memory:  MachineParams.Memory,
		Key:     sql.NullString{String: MachineParams.Key, Valid: true},
		OwnerID: ownerID,
	}

	machine, err := database.Queries.CreateMachine(req.Context(), createParams)
	if err != nil {
		http.Error(res, "Failed to create machine", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(res).Encode(machine)
}

func ListOfFreeMachines(res http.ResponseWriter, req *http.Request) {
	machines, err := database.Queries.ListAvailableMachines(req.Context())
	if err != nil {
		http.Error(res, "Failed to get machines", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(res).Encode(machines)
}

func GetMachineByID(res http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	if id == "" {
		http.Error(res, "ID is required", http.StatusBadRequest)
		return
	}

	machineID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(res, "Invalid ID", http.StatusBadRequest)
		return
	}

	machine, err := database.Queries.GetMachineByID(req.Context(), int32(machineID))
	if err != nil {
		http.Error(res, "Failed to get machine", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(res).Encode(machine)
}

func BuyMachine(res http.ResponseWriter, req *http.Request) {

}
