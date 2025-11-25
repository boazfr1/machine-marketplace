package machine

import (
	json "encoding/json"
	fmt "fmt"
	http "net/http"
	sync "sync"

	websocket "github.com/gorilla/websocket"
	ssh "golang.org/x/crypto/ssh"
)

type Connection struct {
	ws        *websocket.Conn
	sshClient *ssh.Client
	mutex     sync.Mutex
}

func NewConnection(ws *websocket.Conn, client *ssh.Client) *Connection {
	return &Connection{
		ws:        ws,
		sshClient: client,
	}
}

type ConnectionManager struct {
	connections map[string]*Connection
	mutex       sync.RWMutex
}

func (m *ConnectionManager) AddConnection(host string, conn *Connection) {
	m.mutex.Lock()
	m.connections[host] = conn
	m.mutex.Unlock()
}

var Manager = &ConnectionManager{
	connections: make(map[string]*Connection),
}

// type chosenMachineParams struct {
// 	Key     string `json:"key"`
// 	Host    string `json:"host"`
// 	SshUser string `json:"ssh_user"`
// 	Type    string `json:"type"`
// }

// type machineName struct {
// 	MachineName string `json:"machine_name"`
// 	OwnerName   string `json:"owner_name"`
// }

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	EnableCompression: true,
}

type returnedMassage struct {
	Massage  string `json:"massage"`
	Location string `json:"Location"`
}

func HandleConnection(host string, conn *Connection) {
	defer func() {
		conn.ws.Close()
		conn.sshClient.Close()

		Manager.mutex.Lock()
		delete(Manager.connections, host)
		Manager.mutex.Unlock()
	}()

	for {
		_, command, err := conn.ws.ReadMessage()
		if err != nil {
			fmt.Printf("Error reading WebSocket message: %v\n", err)
			return
		}

		output, err := ExecuteCommand(conn, string(command))
		if err != nil {
			conn.ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Error: %v", err)))
			continue
		}

		pwd, err := ExecuteCommand(conn, string("pwd"))
		if err != nil {
			conn.ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Error: %v", err)))
			continue
		}

		massage := returnedMassage{
			Massage:  output,
			Location: pwd,
		}

		jsonData, err := json.Marshal(massage)
		if err != nil {
			fmt.Printf("Error marshaling to JSON: %v\n", err)
			return
		}

		if err := conn.ws.WriteMessage(websocket.TextMessage, jsonData); err != nil {
			fmt.Printf("Error writing response: %v\n", err)
			return
		}
	}
}
