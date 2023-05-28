package chat

import (
	"net/http"
	"time"

	"github.com/gobwas/ws"
	"github.com/gocql/gocql"
)

// chat to store send and recieve messages
var hub *Hub

func init() {
	//  create keyspace and tables

	hub = NewHub()
	go hub.run()
}

type ChatMessage struct {
	ConversationId gocql.UUID `json:"conversation_id" cql:"conversation_id"`
	MessageId      gocql.UUID `json:"messaage_id" cql:"message_id"`
	SenderId       gocql.UUID `json:"sender_id" cql:"sender_id"`
	RecipientId    gocql.UUID `json:"recipient_id" cql:"recipient_id"`
	Message        string     `json:"message" cql:"message"`
	Timestamp      time.Time  `json:"timestamp" cql:"timestamp"`
}

func (s *ChatMessage) ModelType() string {
	return "ChatMessage"
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
