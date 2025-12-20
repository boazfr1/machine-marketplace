package order

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"strconv"
	"time"

	db "machine-marketplace/internal/DB/generated"
	"machine-marketplace/internal/machine"
	"machine-marketplace/pkg/kafka"
)

var (
	serviceLogger = slog.New(slog.NewTextHandler(os.Stdout, nil))
)

func (m *Module) ListOfFreeMachinesService(ctx context.Context, cpuStr string, ramStr string, gpuStr string) ([]db.Machine, error) {
	if cpuStr != "" || ramStr != "" || gpuStr != "" {
		serviceLogger.Info("ListOfFreeMachinesService - filtering machines", "cpu", cpuStr, "ram", ramStr, "gpu", gpuStr)

		cpu, err := convertToInt(cpuStr)
		if err != nil {
			serviceLogger.Error("ListOfFreeMachinesService - failed to convert cpu", "error", err, "cpu", cpuStr)
			return nil, err
		}
		ram, err := convertToInt(ramStr)
		if err != nil {
			serviceLogger.Error("ListOfFreeMachinesService - failed to convert ram", "error", err, "ram", ramStr)
			return nil, err
		}
		gpu, err := convertToInt(gpuStr)
		if err != nil {
			serviceLogger.Error("ListOfFreeMachinesService - failed to convert gpu", "error", err, "gpu", gpuStr)
			return nil, err
		}

		machines, err := m.DB.FilterAvailableMachines(ctx, db.FilterAvailableMachinesParams{
			Column1: cpu,
			Column2: ram,
			Column3: gpu,
		})
		if err != nil {
			serviceLogger.Error("ListOfFreeMachinesService - failed to filter machines", "error", err)
			return nil, err
		}
		serviceLogger.Info("ListOfFreeMachinesService - filtered machines retrieved", "count", len(machines))
		return machines, nil
	}

	serviceLogger.Info("ListOfFreeMachinesService - listing all available machines")
	machines, err := m.DB.ListAvailableMachines(ctx)
	if err != nil {
		serviceLogger.Error("ListOfFreeMachinesService - failed to list machines", "error", err)
		return nil, err
	}
	serviceLogger.Info("ListOfFreeMachinesService - all available machines retrieved", "count", len(machines))
	return machines, nil
}

func (m *Module) CreateMachineService(ctx context.Context, params machine.MachineParams, ownerID int32) (db.Machine, error) {
	serviceLogger.Info("CreateMachineService - validating SSH connection", "host", params.Host, "user", params.SshUser)

	err := machine.TriedToConnectForFirstTime(params.Host, params.SshUser, params.Key)
	if err != nil {
		serviceLogger.Error("CreateMachineService - SSH validation failed", "error", err, "host", params.Host)
		return db.Machine{}, err
	}

	serviceLogger.Info("CreateMachineService - SSH validation successful, creating machine", "name", params.Name, "owner_id", ownerID)

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

	newMachine, err := m.DB.CreateMachine(ctx, createParams)
	if err != nil {
		serviceLogger.Error("CreateMachineService - failed to create machine", "error", err, "name", params.Name)
		return db.Machine{}, err
	}

	serviceLogger.Info("CreateMachineService - machine created successfully", "machine_id", newMachine.ID, "name", newMachine.Name)
	return newMachine, nil
}

func (m *Module) GetOwnedMachinesService(ctx context.Context, ownerID int32) ([]db.Machine, error) {
	serviceLogger.Info("GetOwnedMachinesService - fetching owned machines", "owner_id", ownerID)

	machines, err := m.DB.ListMachinesByOwnerID(ctx, ownerID)
	if err != nil {
		serviceLogger.Error("GetOwnedMachinesService - failed to fetch owned machines", "error", err, "owner_id", ownerID)
		return nil, err
	}

	serviceLogger.Info("GetOwnedMachinesService - owned machines retrieved", "owner_id", ownerID, "count", len(machines))
	return machines, nil
}

func (m *Module) GetBoughtMachinesService(ctx context.Context, buyerID int32) ([]db.ListMachinesByBuyerIDRow, error) {
	serviceLogger.Info("GetBoughtMachinesService - fetching bought machines", "buyer_id", buyerID)

	buyerIDParam := sql.NullInt32{
		Int32: buyerID,
		Valid: true,
	}

	machines, err := m.DB.ListMachinesByBuyerID(ctx, buyerIDParam)
	if err != nil {
		serviceLogger.Error("GetBoughtMachinesService - failed to fetch bought machines", "error", err, "buyer_id", buyerID)
		return nil, err
	}

	serviceLogger.Info("GetBoughtMachinesService - bought machines retrieved", "buyer_id", buyerID, "count", len(machines))
	return machines, nil
}

type BuyMachineResult struct {
	Machine        db.Machine `json:"machine"`
	DealExpiration time.Time  `json:"deal_expiration"`
}

func (m *Module) BuyMachineService(ctx context.Context, machineID int32, buyerID int32, dealDurationHours int) (BuyMachineResult, error) {
	serviceLogger.Info("BuyMachineService - starting purchase process", "machine_id", machineID, "buyer_id", buyerID, "duration_hours", dealDurationHours)

	machineDetails, err := m.DB.GetMachineByID(ctx, machineID)
	if err != nil {
		serviceLogger.Error("BuyMachineService - failed to get machine details", "error", err, "machine_id", machineID)
		return BuyMachineResult{}, err
	}

	if machineDetails.BuyerID.Valid {
		serviceLogger.Error("BuyMachineService - machine already purchased", "machine_id", machineID, "existing_buyer_id", machineDetails.BuyerID.Int32)
		return BuyMachineResult{}, sql.ErrNoRows
	}

	// Send purchase event to Kafka - another service will handle the actual machine update
	dealExpiration := time.Now().Add(time.Duration(dealDurationHours) * time.Hour)
	purchaseEvent := kafka.PurchaseEvent{
		MachineID:      machineID,
		BuyerID:        buyerID,
		DealExpiration: dealExpiration,
		PurchaseTime:   time.Now(),
		MachineName:    machineDetails.Name,
	}

	serviceLogger.Info("BuyMachineService - sending purchase event to Kafka", "machine_id", machineID, "buyer_id", buyerID, "expiration", dealExpiration)

	if err := m.KafkaProducer.SendPurchaseEvent(ctx, purchaseEvent); err != nil {
		serviceLogger.Error("BuyMachineService - failed to send purchase event", "error", err, "machine_id", machineID)
		return BuyMachineResult{}, err
	}

	serviceLogger.Info("BuyMachineService - purchase completed successfully", "machine_id", machineID, "buyer_id", buyerID)
	return BuyMachineResult{
		Machine:        machineDetails,
		DealExpiration: dealExpiration,
	}, nil
}

func convertToInt(str string) (int32, error) {
	if str == "" {
		return 0, nil
	}
	intVal, err := strconv.Atoi(str)
	if err != nil {
		serviceLogger.Error("convertToInt - failed to convert string to int", "error", err, "value", str)
		return 0, err
	}
	return int32(intVal), nil
}
