// BoughtMachines.tsx
import { useEffect, useState } from "react";
import Machine from "../components/Machine";
import "./style/Feed.css";
import Sidebar from "../../dashboard/SideBar";
import api from "../../global/api";
import { Dialog } from "@mui/material";
import Terminal from "../components/Terminal";
import { MachineType } from "../../global/types";



const BoughtMachines = () => {
   const [availableMachine, setAvailableMachine] = useState<MachineType[]>([]);
   const [isLoading, setIsLoading] = useState(true);
   const [isTerminalOpen, setIsTerminalOpen] = useState(false);
   const [selectedMachine, setSelectedMachine] = useState<MachineType | null>(null);

   const getAllAvailableMachine = async () => {
       try {
           const { data } = await api<MachineType[]>('/api/v1/machine/bought-machines');
           console.log("data = ", data);
           setAvailableMachine(data);
           console.log("availableMachine = ", availableMachine);
       } catch {
           throw new Error("Failed to get bought machines");
       }
       finally {
           setIsLoading(false);
       }
   }

   useEffect(() => {
       getAllAvailableMachine();
   }, []);

   const handleOpenTerminal = (machine: MachineType) => {
       setSelectedMachine(machine);
       setIsTerminalOpen(true);
   };

   const handleCloseTerminal = () => {
       setIsTerminalOpen(false);
       setSelectedMachine(null);
   };

   return (
       <div className="feed-page">
           <Sidebar />
           {availableMachine ?
               <div className={`feed-container ${isLoading ? 'loading' : ''}`}>
                   {!isLoading && availableMachine.map((machine, index) => (
                       <div key={index} className="machine-card">
                           <Machine
                               Name={machine.Name}
                               Ram={machine.Ram}
                               Cpu={machine.Cpu}
                               Memory={machine.Memory}
                               OwnerID={machine.OwnerID}
                           />
                           <button onClick={() => handleOpenTerminal(machine)}>
                               open terminal
                           </button>
                       </div>
                   ))}
               </div> :
               <div>There is no available machines to show</div>
           }

           <Dialog
               open={isTerminalOpen}
               onClose={handleCloseTerminal}
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
               {selectedMachine && (
                   <Terminal
                       onClose={handleCloseTerminal}
                       machineName={selectedMachine.Name}
                       ownerName={selectedMachine.OwnerID}
                   />
               )}
           </Dialog>
       </div>
   );
};

export default BoughtMachines;