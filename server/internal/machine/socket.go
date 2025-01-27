package machine

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	websocket "github.com/gorilla/websocket"
	ssh "golang.org/x/crypto/ssh"
)

type Connection struct {
	ws        *websocket.Conn
	sshClient *ssh.Client
	// sshSession *ssh.Session
	mutex sync.Mutex
}

type ConnectionManager struct {
	connections map[string]*Connection
	mutex       sync.RWMutex
}

var Manager = &ConnectionManager{
	connections: make(map[string]*Connection),
}

type chosenMachineParams struct {
	Key     string `json:"key"`
	Host    string `json:"host"`
	SshUser string `json:"ssh_user"`
}

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type returnedMassage struct {
	Massage  string `json:"massage"`
	Location string `json:"Location"`
}

func WebSocketHandler(res http.ResponseWriter, req *http.Request) {
	var params chosenMachineParams
	if err := json.NewDecoder(req.Body).Decode(&params); err != nil {
		http.Error(res, "Invalid request body", http.StatusBadRequest)
		return
	}

	if params.Key == "" || params.Host == "" || params.SshUser == "" {
		http.Error(res, "key, host, and ssh_user are required", http.StatusBadRequest)
		return
	}

	wsConn, err := Upgrader.Upgrade(res, req, nil)
	if err != nil {
		fmt.Println("Error upgrading:", err)
		return
	}

	sshClient, err := CreateSSHClient(params.Host, params.SshUser, params.Key)
	if err != nil {
		wsConn.Close()
		fmt.Printf("Error creating SSH client: %v\n", err)
		return
	}

	conn := &Connection{
		ws:        wsConn,
		sshClient: sshClient,
	}

	Manager.mutex.Lock()
	Manager.connections[params.Host] = conn
	Manager.mutex.Unlock()

	go handleConnection(params.Host, conn)
}

func handleConnection(host string, conn *Connection) {
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

// func handleConnection(wsConn *websocket.Conn, key string, user string, host string) {
// 	for {
// 		_, massage, err := wsConn.ReadMessage()
// 		if err != nil {
// 			fmt.Println("Error reading message:", err)
// 			break
// 		}

// 		if err := wsConn.WriteMessage(websocket.TextMessage, massage); err != nil {
// 			fmt.Println("Error writing message:", err)
// 			break
// 		}

// 	}
// }
