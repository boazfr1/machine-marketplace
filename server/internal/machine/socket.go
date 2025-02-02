package machine

import (
	json "encoding/json"
	fmt "fmt"
	db "machine-marketplace/internal/DB/generated"
	middleware "machine-marketplace/internal/middleware"
	database "machine-marketplace/pkg/database"
	http "net/http"
	"strconv"
	sync "sync"

	jwt "github.com/dgrijalva/jwt-go/v4"
	websocket "github.com/gorilla/websocket"
	ssh "golang.org/x/crypto/ssh"
)

type Connection struct {
	ws        *websocket.Conn
	sshClient *ssh.Client
	mutex     sync.Mutex
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
	Type    string `json:"type"`
}

type machineName struct {
	MachineName string `json:"machine_name"`
	OwnerName   string `json:"owner_name"`
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

	claims := req.Context().Value(middleware.ClaimsContextKey).(*jwt.StandardClaims)

	query := req.URL.Query()

	createParams := db.GetMachineByNameAndOwnerParams{
		Name:   query.Get("machine_name"),
		Name_2: query.Get("owner_name"),
	}

	params, err := database.Queries.GetMachineByNameAndOwner(req.Context(), createParams)
	if err != nil {
		http.Error(res, "Machine not found", http.StatusNotFound)
		return
	}

	num, err := strconv.Atoi(claims.Issuer)
	if err != nil {
		http.Error(res, "Unauthorized", http.StatusUnauthorized)
		return
	}
	buyerID := int32(num)

	if params.BuyerID.Int32 != buyerID {
		http.Error(res, "You are not the owner of this machine", http.StatusForbidden)
		return
	}

	wsConn, err := Upgrader.Upgrade(res, req, nil)
	if err != nil {
		fmt.Println("Error upgrading:", err)
		return
	}

	sshClient, err := CreateSSHClient(params.Host, params.SshUser, params.Key.String)
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
