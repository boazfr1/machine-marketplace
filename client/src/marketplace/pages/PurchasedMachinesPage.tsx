import { useEffect, useState } from "react";
import MarketplaceLayout from "../components/MarketplaceLayout";
import { MachineType } from "../../global/types";
import "./PurchasedMachinesPage.css";

interface PurchasedMachine extends MachineType {
  purchaseDate: string;
  status: 'active' | 'expired' | 'pending';
  expiryDate: string;
  totalCost: number;
}

const PurchasedMachinesPage = () => {
  const [purchasedMachines, setPurchasedMachines] = useState<PurchasedMachine[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const fetchPurchasedMachines = async () => {
      try {
        await new Promise(resolve => setTimeout(resolve, 1000));
        
        // Mock data
        const mockPurchases: PurchasedMachine[] = [
          {
            Name: "AI Training Server",
            Ram: 64,
            Cpu: 32,
            Memory: 2000,
            OwnerID: 456,
            purchaseDate: "2024-01-15",
            status: 'active',
            expiryDate: "2024-02-15",
            totalCost: 156.50
          },
          {
            Name: "Data Processing Cluster",
            Ram: 128,
            Cpu: 48,
            Memory: 4000,
            OwnerID: 789,
            purchaseDate: "2024-01-10",
            status: 'expired',
            expiryDate: "2024-01-20",
            totalCost: 89.25
          }
        ];
        
        setPurchasedMachines(mockPurchases);
      } catch (error) {
        console.error("Failed to fetch purchased machines:", error);
      } finally {
        setIsLoading(false);
      }
    };

    fetchPurchasedMachines();
  }, []);

  const getStatusBadge = (status: string) => {
    const statusClasses = {
      active: 'status-active',
      expired: 'status-expired',
      pending: 'status-pending'
    };
    
    return (
      <span className={`status-badge ${statusClasses[status as keyof typeof statusClasses]}`}>
        {status.charAt(0).toUpperCase() + status.slice(1)}
      </span>
    );
  };

  const handleConnect = (machine: PurchasedMachine) => {
    console.log("Connect to machine:", machine.Name);
    // Add connection logic here
  };

  const handleExtend = (machine: PurchasedMachine) => {
    console.log("Extend rental for:", machine.Name);
    // Add extend logic here
  };

  return (
    <MarketplaceLayout>
      <div className="purchased-machines-page">
        <div className="page-header">
          <div className="header-content">
            <h1 className="page-title">My Rentals</h1>
            <p className="page-subtitle">
              Manage your rented machines and access your computing resources
            </p>
          </div>
          
          <div className="header-stats">
            <div className="stat-card">
              <span className="stat-value">
                {purchasedMachines.filter(m => m.status === 'active').length}
              </span>
              <span className="stat-label">Active Rentals</span>
            </div>
            <div className="stat-card">
              <span className="stat-value">
                ${purchasedMachines.reduce((sum, m) => sum + m.totalCost, 0).toFixed(2)}
              </span>
              <span className="stat-label">Total Spent</span>
            </div>
          </div>
        </div>

        <div className="page-content">
          {isLoading ? (
            <div className="loading-state">
              <div className="loading-spinner-large"></div>
              <p>Loading your rentals...</p>
            </div>
          ) : purchasedMachines.length > 0 ? (
            <div className="rentals-list">
              {purchasedMachines.map((machine, index) => (
                <div key={`${machine.Name}-${index}`} className="rental-card">
                  <div className="rental-header">
                    <div className="rental-info">
                      <h3 className="rental-name">{machine.Name}</h3>
                      <div className="rental-meta">
                        <span className="owner-info">Provider #{machine.OwnerID}</span>
                        {getStatusBadge(machine.status)}
                      </div>
                    </div>
                    <div className="rental-cost">
                      <span className="cost-amount">${machine.totalCost}</span>
                      <span className="cost-label">Total Cost</span>
                    </div>
                  </div>

                  <div className="rental-specs">
                    <div className="spec-item">
                      <span className="spec-icon">🧠</span>
                      <span className="spec-text">{machine.Cpu} CPU Cores</span>
                    </div>
                    <div className="spec-item">
                      <span className="spec-icon">💾</span>
                      <span className="spec-text">{machine.Ram} GB RAM</span>
                    </div>
                    <div className="spec-item">
                      <span className="spec-icon">💿</span>
                      <span className="spec-text">{machine.Memory} GB Storage</span>
                    </div>
                  </div>

                  <div className="rental-timeline">
                    <div className="timeline-item">
                      <span className="timeline-label">Started:</span>
                      <span className="timeline-date">
                        {new Date(machine.purchaseDate).toLocaleDateString()}
                      </span>
                    </div>
                    <div className="timeline-item">
                      <span className="timeline-label">
                        {machine.status === 'expired' ? 'Expired:' : 'Expires:'}
                      </span>
                      <span className="timeline-date">
                        {new Date(machine.expiryDate).toLocaleDateString()}
                      </span>
                    </div>
                  </div>

                  <div className="rental-actions">
                    {machine.status === 'active' && (
                      <>
                        <button 
                          onClick={() => handleConnect(machine)}
                          className="btn btn-primary"
                        >
                          Connect
                        </button>
                        <button 
                          onClick={() => handleExtend(machine)}
                          className="btn btn-outline"
                        >
                          Extend Rental
                        </button>
                      </>
                    )}
                    {machine.status === 'expired' && (
                      <button className="btn btn-secondary" disabled>
                        Rental Expired
                      </button>
                    )}
                    {machine.status === 'pending' && (
                      <button className="btn btn-outline" disabled>
                        Setting Up...
                      </button>
                    )}
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <div className="empty-state">
              <div className="empty-icon">💰</div>
              <h3 className="empty-title">No rentals yet</h3>
              <p className="empty-description">
                You haven't rented any machines yet. Browse the marketplace to find the perfect machine for your needs.
              </p>
              <button className="btn btn-primary">
                Browse Marketplace
              </button>
            </div>
          )}
        </div>
      </div>
    </MarketplaceLayout>
  );
};

export default PurchasedMachinesPage;
