import { FC } from "react";
import "./MachineCard.css";

interface MachineCardProps {
    Ram: number;
    Cpu: number;
    Gpu: number;
    Memory: number;
    Name: string;
    OwnerID: number;
    onClick?: () => void;
    showActions?: boolean;
    price?: number;
}

const MachineCard: FC<MachineCardProps> = ({
    Name,
    OwnerID,
    Ram,
    Cpu,
    Gpu,
    Memory,
    onClick,
    showActions = true,
    price = 0.25
}) => {
    const handleCardClick = () => {
        if (onClick) {
            onClick();
        }
    };

    const handleActionClick = (e: React.MouseEvent, action: string) => {
        e.stopPropagation();
        console.log(`${action} clicked for machine: ${Name}`);
    };

    return (
        <div className="machine-card" onClick={handleCardClick}>
            <div className="machine-card-header">
                <div className="machine-info">
                    <h3 className="machine-name">{Name}</h3>
                    <div className="machine-owner">
                        <span className="owner-label">Owner:</span>
                        <span className="owner-id">#{OwnerID}</span>
                    </div>
                </div>
                <div className="machine-status">
                    <span className="status-badge status-available">Available</span>
                </div>
            </div>

            <div className="machine-specs">
                <div className="spec-grid">
                    <div className="spec-item">
                        <div className="spec-icon">🧠</div>
                        <div className="spec-details">
                            <span className="spec-label">CPU Cores</span>
                            <span className="spec-value">{Cpu}</span>
                        </div>
                    </div>

                    <div className="spec-item">
                        <div className="spec-icon">💾</div>
                        <div className="spec-details">
                            <span className="spec-label">RAM</span>
                            <span className="spec-value">{Ram} GB</span>
                        </div>
                    </div>

                    <div className="spec-item">
                        <div className="spec-icon">🎮</div>
                        <div className="spec-details">
                            <span className="spec-label">GPU</span>
                            <span className="spec-value">{Gpu > 0 ? `${Gpu} GPU` : 'No GPU'}</span>
                        </div>
                    </div>

                    <div className="spec-item">
                        <div className="spec-icon">💿</div>
                        <div className="spec-details">
                            <span className="spec-label">Storage</span>
                            <span className="spec-value">{Memory} GB</span>
                        </div>
                    </div>
                </div>
            </div>

            <div className="machine-footer">
                <div className="machine-pricing">
                    <span className="price-label">Starting at</span>
                    <span className="price-value">${price}/hour</span>
                </div>

                {showActions && (
                    <div className="machine-actions">
                        <button
                            className="btn btn-outline btn-sm"
                            onClick={(e) => handleActionClick(e, 'view')}
                        >
                            View Details
                        </button>
                        <button
                            className="btn btn-primary btn-sm"
                            onClick={(e) => handleActionClick(e, 'rent')}
                        >
                            Rent Now
                        </button>
                    </div>
                )}
            </div>
        </div>
    );
};

export default MachineCard;
