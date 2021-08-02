package socket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var Connection *websocket.Conn

func Handle(w http.ResponseWriter, r *http.Request) {
	Connection, _ = upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity

	for {
		// Read message from browser
		msgType, msg, err := Connection.ReadMessage()
		if err != nil {
			return
		}

		// Print the message to the console
		fmt.Printf("%s sent: %s\n", Connection.RemoteAddr(), string(msg))

		// Write message back to browser
		if err = Connection.WriteMessage(msgType, msg); err != nil {
			return
		}
	}
}
