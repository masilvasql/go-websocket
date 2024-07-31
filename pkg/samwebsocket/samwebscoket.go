package samwebsocket

import (
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

var clients = make(map[string][]*websocket.Conn)
var clientsLock = &sync.Mutex{}

type Upgrader struct {
	ReadBufferSize  int
	WriteBufferSize int
	CheckOrigin     func(r *http.Request) bool
}

func NewUpgrader(readBufferSize, writerBufferSize int, checkOrigin bool) *websocket.Upgrader {
	return &websocket.Upgrader{
		ReadBufferSize:  readBufferSize,
		WriteBufferSize: writerBufferSize,
		CheckOrigin: func(r *http.Request) bool {
			return checkOrigin
		},
	}
}

func AddClient(sessionID string, conn *websocket.Conn) {
	clientsLock.Lock()
	defer clientsLock.Unlock()
	clients[sessionID] = append(clients[sessionID], conn)
}

func RemoveClient(sessionID string, conn *websocket.Conn) {
	clientsLock.Lock()
	defer clientsLock.Unlock()
	for i, client := range clients[sessionID] {
		if client == conn {
			clients[sessionID] = append(clients[sessionID][:i], clients[sessionID][i+1:]...)
			break
		}
	}
}

func HandleMessages(sessionID string, conn *websocket.Conn) {
	defer func() {
		RemoveClient(sessionID, conn)
		conn.Close()
	}()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			return
		}

		broadcastMessage(sessionID, messageType, p)
	}
}

func broadcastMessage(sessionID string, messageType int, message []byte) {
	clientsLock.Lock()
	defer clientsLock.Unlock()
	for _, client := range clients[sessionID] {
		if err := client.WriteMessage(messageType, message); err != nil {
			client.Close()
			RemoveClient(sessionID, client)
		}
	}
}
