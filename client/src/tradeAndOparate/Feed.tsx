import axios from "axios";
import { useEffect, useState } from "react";
import Machine from "./Machine";
import "./Feed.css";
import { useNavigate } from "react-router-dom";
import Sidebar from "../dashboard/SideBar";


type MachineType = {
    Name: string;
    RAM: number;
    CPU: number;
    Storage: number;
};

const Feed = () => {
    const [availableMachine, setAvailableMachine] = useState<MachineType[]>([]);
    const [isLoading, setIsLoading] = useState(true);

    const navigate = useNavigate();

    

    const getAllAvailableMachine = async () => {
        try {
            const { data } = await axios<MachineType[]>('http://localhost:3001/api/v1/machine');
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
                            RAM={machine.RAM}
                            CPU={machine.CPU}
                            Storage={machine.Storage}
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