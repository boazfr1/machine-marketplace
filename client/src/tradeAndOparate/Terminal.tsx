import { FC, useState, KeyboardEvent, useEffect } from 'react';
import { IconButton } from '@mui/material';
import CloseIcon from '@mui/icons-material/Close';
import './Terminal.css';


interface TerminalProps {
    onClose: () => void;
    machineName: string;
    sshKey: string;
    host: string;
    sshUser: string;
}

interface CommandEntry {
    command: string;
    output: string;
}

const Terminal: FC<TerminalProps> = ({ onClose, machineName, sshKey, host, sshUser }) => {
    const [commandHistory, setCommandHistory] = useState<CommandEntry[]>([]);
    const [currentCommand, setCurrentCommand] = useState('');
    const [pwd, setPwd] = useState('')
    const [socketCon, setSocketCon] = useState<WebSocket>();

    useEffect(() => {
        connectWebSocket();

        return () => {
            if (socketCon) {
                socketCon.close();
            }
        };
    }, []);

    const sendCommand = () => {
        const msg = {
            command: currentCommand,
        };
        if (socketCon != undefined) {
            socketCon.send(JSON.stringify(msg));
        }
        return ''
    };

    const connectWebSocket = async () => {
        try {
            const response = await fetch('/api/v1/machine/connect', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    key: sshKey,
                    host: host,
                    ssh_user: sshUser
                })
            });

            if (response.ok) {
                const ws = new WebSocket(`ws://${window.location.host}/api/v1/machine/connect`);
                
                ws.onmessage = (event) => {
                    try {
                        const response = JSON.parse(event.data);
                        setCommandHistory(prev => [...prev, {
                            command: currentCommand,
                            output: response.massage
                        }]);
                        setPwd(response.Location);
                    } catch (error) {
                        console.error('Error parsing response:', error);
                    }
                };

                ws.onopen = () => {
                    console.log('WebSocket connected');
                };

                ws.onerror = (error) => {
                    console.error('WebSocket error:', error);
                };

                setSocketCon(ws);
            }
        } catch (error) {
            console.error('Connection error:', error);
        }
    };

    // const processCommand = (command: string): string => {
    //     // Add your command processing logic here
    //     switch (command.toLowerCase()) {
    //         case 'help':
    //             return 'Available commands:\nhelp - Show this help message\nclear - Clear terminal\necho [text] - Echo text back\nexit - Close terminal';
    //         case 'clear':
    //             setCommandHistory([]);
    //             return '';
    //         case 'exit':
    //             onClose();
    //             return '';
    //         default:
    //             if (command.toLowerCase().startsWith('echo ')) {
    //                 return command.slice(5);
    //             }
    //             return `Command not found: ${command}`;
    //     }
    // };

    const handleKeyPress = (e: KeyboardEvent<HTMLInputElement>) => {
        if (e.key === 'Enter' && currentCommand) {
            const output = sendCommand();
            setCommandHistory([...commandHistory, { command: currentCommand, output }]);
            setCurrentCommand('');
        }
    };

    return (
        <div className="terminal-container">
            <div className="terminal-header">
                <span>{machineName} - Terminal</span>
                <IconButton
                    onClick={onClose}
                    size="small"
                    sx={{ color: '#fff' }}
                >
                    <CloseIcon />
                </IconButton>
            </div>

            <div className="terminal-content">
                {commandHistory.map((entry, index) => (
                    <div key={index}>
                        <div className="command-line">
                            <span className="prompt">$ </span>
                            <span>{entry.command}</span>
                        </div>
                        {entry.output && (
                            <div className="command-output">{entry.output}</div>
                        )}
                    </div>
                ))}
                <div className="current-line">
                    <span className="prompt">{pwd}$ </span>
                    <input
                        type="text"
                        value={currentCommand}
                        onChange={(e) => setCurrentCommand(e.target.value)}
                        onKeyPress={handleKeyPress}
                        autoFocus
                        spellCheck={false}
                    />
                </div>
            </div>
        </div>
    );
};

export default Terminal;