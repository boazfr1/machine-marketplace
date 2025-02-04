import { useLocation, useNavigate } from 'react-router-dom';
import './style/MachineDetails.css'
import { MachineType } from '../../global/types';

const MachineDetails = () => {
    const location = useLocation();
    const navigate = useNavigate();
    const machine = location.state?.machine as MachineType;

    if (!machine) {
        navigate('/');
        return null;
    }

    const buyMachine = () => {
        console.log("Buy Machine");
    }

    return (
        <div className="machine-details-container">
            <h1>{machine.Name}</h1>
            <div className="machine-specs">
                <div className="spec-item">
                    <label>RAM:</label>
                    <span>{machine.Ram} GB</span>
                </div>
                <div className="spec-item">
                    <label>CPU Cores:</label>
                    <span>{machine.Cpu}</span>
                </div>
                <div className="spec-item">
                    <label>Storage:</label>
                    <span>{machine.Memory} GB</span>
                </div>
                <div className="spec-item">
                    <label>Owner:</label>
                    <span>{machine.OwnerID}</span>
                </div>
            </div>
            
            <div className="button-container">
                <button
                    className="connect-button"
                    onClick={buyMachine}
                >
                    Buy
                </button>
            </div>
        </div>
    );
};

export default MachineDetails;