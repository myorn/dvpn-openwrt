package main

import (
	"embed"
	"github.com/solarlabsteam/dvpn-openwrt/controllers"
	"github.com/solarlabsteam/dvpn-openwrt/services/auth"
	"github.com/solarlabsteam/dvpn-openwrt/services/socket"
	"io/fs"
	"net/http"
	"os"
)

//go:embed public

var public embed.FS

func main() {
	if _, homeSet := os.LookupEnv("HOME"); !homeSet {
		os.Setenv("PATH", "/usr/sbin:/usr/bin:/sbin:/bin:")
		os.Setenv("HOME", "/root")
	}
	// for development: serve static assets from public folder
	//publicFS := http.FileServer(http.Dir("./public"))

	// for production: embed static assets into binary
	publicDir, _ := fs.Sub(public, "public")
	publicFS := http.FileServer(http.FS(publicDir))

	http.Handle("/", auth.BasicAuthForHandler(publicFS)) // serve embedded static assets
	http.HandleFunc("/api/node/start/stream", auth.BasicAuth(controllers.StartNodeStreamStd))
	http.HandleFunc("/api/node", auth.BasicAuth(controllers.GetNode))
	http.HandleFunc("/api/node/kill", auth.BasicAuth(controllers.KillNode))
	http.HandleFunc("/api/config", auth.BasicAuth(controllers.Config))
	http.HandleFunc("/api/socket", auth.BasicAuth(socket.Handle))
	http.HandleFunc("/api/keys", auth.BasicAuth(controllers.ListKeys))
	http.HandleFunc("/api/keys/add", auth.BasicAuth(controllers.AddRecoverKeys))

	if err := http.ListenAndServe(":9000", nil); err != nil {
		panic("failed to start server")
	}
}
