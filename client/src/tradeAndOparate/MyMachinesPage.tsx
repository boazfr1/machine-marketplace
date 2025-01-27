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

const MyMachinesPage = () => {
    const [availableMachine, setAvailableMachine] = useState<MachineType[]>([]);
    const [isLoading, setIsLoading] = useState(true);

    const navigate = useNavigate();

    const mock = [
        {
            Name: "Ubuntu Server 1",
            RAM: 16,
            CPU: 4,
            Storage: 512
        },
        {
            Name: "Development Machine",
            RAM: 32,
            CPU: 8,
            Storage: 1024
        },
        {
            Name: "Testing Environment",
            RAM: 8,
            CPU: 2,
            Storage: 256
        },
        {
            Name: "Production Server",
            RAM: 64,
            CPU: 16,
            Storage: 2048
        },
        {
            Name: "Database Server",
            RAM: 128,
            CPU: 32,
            Storage: 4096
        },
        {
            Name: "Staging Server",
            RAM: 16,
            CPU: 4,
            Storage: 512
        }
    ];

    const getAllAvailableMachine = async () => {
        try {
            const { data } = await axios<MachineType[]>('http://localhost:3001/api/v1/machine/my-machines');
            if (Array.isArray(data)) {
                setAvailableMachine(data);
            } else {
                setAvailableMachine(mock);
            }
        } catch {
            setAvailableMachine(mock)
        }
        finally {
            setIsLoading(false);
        }
    }

    useEffect(() => {

        getAllAvailableMachine();
    }, []);

    const navigateToMachinePage = (machine: MachineType) => {        
        // navigate('/machine', {
        //     state: {
        //         machine: machine
        //     }
        // });
    }

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

export default MyMachinesPage;