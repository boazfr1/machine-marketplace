export interface Machine {
    Name: string;
    Ram: number;
    Cpu: number;
    Memory: number;
    OwnerID: number;
}

// Keep the old type for backward compatibility during transition
export type MachineType = Machine;