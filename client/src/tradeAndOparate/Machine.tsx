// Machine.tsx
import { FC } from "react";
import "./Machine.css"

interface MachineProps {
  Ram: number;
  Cpu: number;
  Memory: number;
  Name: string;
  OwnerID: number;
}

const Machine: FC<MachineProps> = ({ Name, OwnerID, Ram, Cpu, Memory }) => {
  return (
    <div className="machine-card">
      <div className="machine-header">
        <h2>{Name}</h2>
        <p>{OwnerID}</p>
      </div>

      <div className="machine-content">
        <div className="specs-grid">
          <div className="spec-item">
            <div>
              <p className="spec-label">RAM</p>
              <p className="spec-value">{Ram} GB</p>
            </div>
          </div>

          <div className="spec-item">
            <div>
              <p className="spec-label">CPU Cores</p>
              <p className="spec-value">{Cpu}</p>
            </div>
          </div>

          <div className="spec-item">
            <div>
              <p className="spec-label">Storage</p>
              <p className="spec-value">{Memory} GB</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Machine;