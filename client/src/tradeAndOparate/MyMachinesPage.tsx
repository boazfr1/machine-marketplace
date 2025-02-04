import { useEffect, useState } from "react";
import Machine from "./Machine";
import "./Feed.css";
import Sidebar from "../dashboard/SideBar";
import api from "../api";
import { useNavigate } from "react-router-dom";
import { MachineType } from "../global/types";


const MyMachinesPage = () => {
    const [availableMachine, setAvailableMachine] = useState<MachineType[]>([]);
    const [isLoading, setIsLoading] = useState(true);

    const navigate = useNavigate();

    const getAllAvailableMachine = async () => {
        try {
            const { data } = await api<MachineType[]>('/api/v1/machine/owned-machines');
                setAvailableMachine(data);
            
        } catch {
            throw new Error("Failed to get owned machines");
        }
        finally {
            setIsLoading(false);
        }
    }

    const navigateToMachinePage = (machine: MachineType) => {
        navigate('/machine', {
            state: {
                machine: machine
            }
        });
    }

    useEffect(() => {
        getAllAvailableMachine();
    }, []);

    return (
        <div className="feed-page">
            <Sidebar />
            {availableMachine ?
                <div className={`feed-container ${isLoading ? 'loading' : ''}`}>
                    {!isLoading && availableMachine.map((machine, index) => (
                        <div
                            key={index}
                            onClick={() => navigateToMachinePage(machine)}
                        >
                            <Machine
                                Name={machine.Name}
                                Ram={machine.Ram}
                                Cpu={machine.Cpu}
                                Memory={machine.Memory}
                                OwnerID={machine.OwnerID}
                            />
                        </div>
                    ))}
                </div> :
                <div>
                    There is no available machines to show
                </div>
            }
            <div>

            </div>
        </div>

    );
};

export default MyMachinesPage;