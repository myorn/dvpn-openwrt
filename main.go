package main

import (
	"github.com/audi70r/dvpn-openwrt/controllers"
	"github.com/audi70r/dvpn-openwrt/services/socket"
	"io"
	"net/http"
)

var DVPNOut io.ReadCloser

func main() {
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)
	http.HandleFunc("/api/node/start", controllers.StartNode)
	http.HandleFunc("/api/node/start/stream", controllers.StartNodeStreamStd)
	http.HandleFunc("/api/node", controllers.GetNode)
	http.HandleFunc("/api/config", controllers.Config)
	http.HandleFunc("/api/socket", socket.Handle)

	http.ListenAndServe(":9000", nil)
}
