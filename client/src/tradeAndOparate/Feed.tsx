import { useEffect, useState } from "react";
import Machine from "./Machine";
import "./Feed.css";
import { useNavigate } from "react-router-dom";
import Sidebar from "../dashboard/SideBar";
import api from "../api";
import { MachineType } from "../global/types";


const Feed = () => {
    const [availableMachine, setAvailableMachine] = useState<MachineType[]>([]);
    const [isLoading, setIsLoading] = useState(true);

    const navigate = useNavigate();

    

    const getAllAvailableMachine = async () => {
        try {
            const { data } = await api<MachineType[]>('/api/v1/machine');
            console.log("data = ", data);
            setAvailableMachine(data);
        } finally {
            setIsLoading(false);
        }
    }

    useEffect(() => {
        getAllAvailableMachine();
    }, []);

    const navigateToMachinePage = (machine: MachineType) => {
        navigate('/machine', {
            state: {
                machine: machine
            }
        });
    }

    return (
        <div className="feed-page">
            <Sidebar/>
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

export default Feed;