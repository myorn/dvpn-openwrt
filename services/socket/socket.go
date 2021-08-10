package socket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Connection struct {
	Socket *websocket.Conn
	mu     sync.Mutex
}

var Conn Connection

func (c *Connection) Send(msg []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	Conn.Socket.WriteMessage(1, msg)
}

func Handle(w http.ResponseWriter, r *http.Request) {
	Conn.Socket, _ = upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity

	for {
		// Read message from browser
		msgType, msg, err := Conn.Socket.ReadMessage()
		if err != nil {
			return
		}

		// Print the message to the console
		fmt.Printf("%s sent: %s\n", Conn.Socket.RemoteAddr(), string(msg))

		// Write message back to browser
		if err = Conn.Socket.WriteMessage(msgType, msg); err != nil {
			return
		}
	}
}
