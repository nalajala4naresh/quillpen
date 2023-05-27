package chat

import (
	"net/http"

	"github.com/gobwas/ws"
)

// chat to store send and recieve messages
var hub *Hub

func init() {
	hub = NewHub()
	go hub.run()
}

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		// Unable to upgrade the request
		w.WriteHeader(http.StatusUpgradeRequired)
		return

	}

	client := NewClient(conn)
	hub.register <- client

	go client.write(hub)
	client.read(hub)
}
