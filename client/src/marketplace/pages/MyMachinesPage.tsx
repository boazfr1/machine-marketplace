import { useEffect, useState } from "react";
import MarketplaceLayout from "../components/MarketplaceLayout";
import MachineCard from "../components/MachineCard";
import { MachineType } from "../../global/types";
import "./MyMachinesPage.css";

const MyMachinesPage = () => {
  const [myMachines, setMyMachines] = useState<MachineType[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    // Simulate API call to get user's machines
    const fetchMyMachines = async () => {
      try {
        // This would be replaced with actual API call
        await new Promise(resolve => setTimeout(resolve, 1000));
        
        // Mock data - replace with actual API call
        const mockMachines: MachineType[] = [
          { Name: "My Development Server", Ram: 32, Cpu: 16, Memory: 1000, OwnerID: 123 },
          { Name: "ML Training Rig", Ram: 64, Cpu: 24, Memory: 2000, OwnerID: 123 },
        ];
        
        setMyMachines(mockMachines);
      } catch (error) {
        console.error("Failed to fetch my machines:", error);
      } finally {
        setIsLoading(false);
      }
    };

    fetchMyMachines();
  }, []);

  const handleEditMachine = (machine: MachineType) => {
    console.log("Edit machine:", machine.Name);
    // Add edit logic here
  };

  const handleDeleteMachine = (machine: MachineType) => {
    console.log("Delete machine:", machine.Name);
    // Add delete logic here
  };

  return (
    <MarketplaceLayout>
      <div className="my-machines-page">
        <div className="page-header">
          <div className="header-content">
            <h1 className="page-title">My Machines</h1>
            <p className="page-subtitle">
              Manage your listed machines and track their performance
            </p>
          </div>
          
          <div className="header-actions">
            <button className="btn btn-primary">
              Add New Machine
            </button>
          </div>
        </div>

        <div className="page-content">
          {isLoading ? (
            <div className="loading-state">
              <div className="loading-spinner-large"></div>
              <p>Loading your machines...</p>
            </div>
          ) : myMachines.length > 0 ? (
            <div className="machines-grid">
              {myMachines.map((machine, index) => (
                <div key={`${machine.Name}-${index}`} className="machine-card-wrapper">
                  <MachineCard
                    Name={machine.Name}
                    Ram={machine.Ram}
                    Cpu={machine.Cpu}
                    Memory={machine.Memory}
                    OwnerID={machine.OwnerID}
                    showActions={false}
                  />
                  <div className="machine-owner-actions">
                    <button 
                      onClick={() => handleEditMachine(machine)}
                      className="btn btn-outline btn-sm"
                    >
                      Edit
                    </button>
                    <button 
                      onClick={() => handleDeleteMachine(machine)}
                      className="btn btn-secondary btn-sm"
                    >
                      Remove
                    </button>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <div className="empty-state">
              <div className="empty-icon">🖥️</div>
              <h3 className="empty-title">No machines listed</h3>
              <p className="empty-description">
                You haven't listed any machines yet. Start by adding your first machine to the marketplace.
              </p>
              <button className="btn btn-primary">
                Add Your First Machine
              </button>
            </div>
          )}
        </div>
      </div>
    </MarketplaceLayout>
  );
};

export default MyMachinesPage;
