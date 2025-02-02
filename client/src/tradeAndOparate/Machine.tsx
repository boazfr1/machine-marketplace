import { FC, useState } from "react";
import { Dialog } from "@mui/material";
import "./Machine.css"
import Terminal from "./Terminal";


interface MachineProps {
  Name: string;
  Owner: string;
  RAM: number;
  CPU: number;
  Storage: number;
}

const Machine: FC<MachineProps> = ({ Name, Owner, RAM, CPU, Storage }) => {
  const [isTerminalOpen, setIsTerminalOpen] = useState(false);

  const ConnectToTerminal = () => {
    setIsTerminalOpen(true);
  }

  const handleClose = () => {
    setIsTerminalOpen(false);
  };



  return (
    <div className="machine-card">
      <div className="machine-header">
        <h2>{Name}</h2>
        <p>{Owner}</p>
      </div>

      <div className="machine-content">
        <div className="specs-grid">
          <div className="spec-item">
            <div>
              <p className="spec-label">RAM</p>
              <p className="spec-value">{RAM} GB</p>
            </div>
          </div>

          <div className="spec-item">
            <div>
              <p className="spec-label">CPU Cores</p>
              <p className="spec-value">{CPU}</p>
            </div>
          </div>

          <div className="spec-item">
            <div>
              <p className="spec-label">Storage</p>
              <p className="spec-value">{Storage}</p>
            </div>
          </div>
        </div>

        <div className="button-container">
          <button
            className="connect-button"
            onClick={ConnectToTerminal}
          >
            Connect to Terminal
          </button>
        </div>
      </div>

      <Dialog
        open={isTerminalOpen}
        onClose={handleClose}
        maxWidth="md"
        fullWidth
        PaperProps={{
          style: {
            backgroundColor: '#1E1E1E',
            color: '#fff',
            borderRadius: '8px'
          }
        }}
      >
        <Terminal
          onClose={handleClose}
          machineName={Name}
          ownerName={Owner}
        />
      </Dialog>


    </div>
  );
};

export default Machine;