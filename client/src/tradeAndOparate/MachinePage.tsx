import { useLocation, useNavigate } from 'react-router-dom';
import './MachinePage.css'

type MachineType = {
    Name: string;
    RAM: number;
    CPU: number;
    Storage: number;
};

const MachinePage = () => {
    const location = useLocation();
    const navigate = useNavigate();
    const machine = location.state?.machine as MachineType;

    if (!machine) {
        navigate('/');
        return null;
    }

    return (
        <div className="machine-details-container">
            <h1>{machine.Name}</h1>
            <div className="machine-specs">
                <div className="spec-item">
                    <label>RAM:</label>
                    <span>{machine.RAM} GB</span>
                </div>
                <div className="spec-item">
                    <label>CPU Cores:</label>
                    <span>{machine.CPU}</span>
                </div>
                <div className="spec-item">
                    <label>Storage:</label>
                    <span>{machine.Storage} GB</span>
                </div>
            </div>
            
            <div className="actions">
                <button onClick={() => {/* Add your action here */}}>
                    Connect to Terminal
                </button>
                <button onClick={() => {/* Add your action here */}}>
                    View Logs
                </button>
                <button onClick={() => navigate('/')}>
                    Back to Feed
                </button>
            </div>
        </div>
    );
};

export default MachinePage;