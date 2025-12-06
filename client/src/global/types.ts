export interface Machine {
    ID: number;
    Name: string;
    Ram: number;
    Cpu: number;
    Gpu: number;
    Memory: number;
    OwnerID: number;
    BuyerID?: number | null;
    Host: string;
    SshUser: string;
}

// Keep the old type for backward compatibility during transition
export type MachineType = Machine;

export interface PurchaseRequest {
    machine_id: number;
    deal_duration_hours: number;
}

export interface PurchaseResponse {
    message: string;
    machine: Machine;
    deal_expiration: string;
}