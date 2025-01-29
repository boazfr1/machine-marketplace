import { FC, useState, KeyboardEvent, useEffect } from 'react';
import { IconButton } from '@mui/material';
import CloseIcon from '@mui/icons-material/Close';
import './Terminal.css';


interface TerminalProps {
    onClose: () => void;
    machineName: string;
}

interface CommandEntry {
    command: string;
    output: string;
}

interface socketResponse {
    response: string;
    location: string;
}

const Terminal: FC<TerminalProps> = ({ onClose, machineName}) => {
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

        const params = {
            machineName
        };
        const queryParams = new URLSearchParams(params).toString();

        const websocket = new WebSocket(`ws://localhost:8080/api/v1/machine/connect?${queryParams}`);
        

        websocket.onopen = () => {
            console.log('Connected to WebSocket');
            websocket.send(JSON.stringify(params));
        };

        websocket.onmessage = (event) => {
            const data: socketResponse = JSON.parse(event.data);
            setPwd(data.location);
            setCommandHistory([...commandHistory, { command: currentCommand, output: data.response }]);
            setCurrentCommand('');
        };

        websocket.onerror = (error) => {
            console.error('WebSocket error:', error);
        };

        setSocketCon(websocket);
    };

    
    const handleKeyPress = (e: KeyboardEvent<HTMLInputElement>) => {
        if (e.key === 'Enter' && currentCommand) {
            sendCommand();
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