package main

import (
	"embed"
	"github.com/audi70r/dvpn-openwrt/controllers"
	"github.com/audi70r/dvpn-openwrt/services/socket"
	"net/http"
)

//go:embed public
var public embed.FS

func main() {
	publicFS := http.FileServer(http.Dir("./public"))

	//publicDir, _ := fs.Sub(public, "public")
	//publicFS := http.FileServer(http.FS(publicDir))

	http.Handle("/", publicFS) // serve embedded static assets
	http.HandleFunc("/api/node/start/stream", controllers.StartNodeStreamStd)
	http.HandleFunc("/api/node", controllers.GetNode)
	http.HandleFunc("/api/node/kill", controllers.KillNode)
	http.HandleFunc("/api/config", controllers.Config)
	http.HandleFunc("/api/socket", socket.Handle)
	http.HandleFunc("/api/keys", controllers.ListKeys)
	http.HandleFunc("/api/keys/add", controllers.AddRecoverKeys)

	if err := http.ListenAndServe(":9000", nil); err != nil {
		panic("failed to start server")
	}
}
