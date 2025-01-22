import { FC } from "react";
import "./Machine.css"


interface MachineProps {
  Name: string;
  RAM: number;
  CPU: number;
  Storage: number;
}

const Machine: FC<MachineProps> = ({Name, RAM, CPU, Storage }) => {

  const ConnectToTerminal = () => {
    console.log("ConnectToTerminal")
  }
  

  return (
    <div className="machine-card">
      <div className="machine-header">
        <h2>{Name}</h2>
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
          
          <button 
            className="action-button"
            onClick={() => {}}
          >
            More Actions
          </button>
        </div>
      </div>

      
    </div>
  );
};

export default Machine;