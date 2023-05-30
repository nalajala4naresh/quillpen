package chat

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gocql/gocql"
	"github.com/gorilla/websocket"

	"github.com/quillpen/accounts"
	"github.com/quillpen/storage"
)

// chat to store send and recieve messages
var hub *Hub

func init() {
	//  create keyspace and tables

	hub = NewHub()
	go hub.run()
}

type Conversation struct {
	ConversationId gocql.UUID `json:"conversation_id" cql:"conversation_id"`
	SenderId       gocql.UUID `json:"sender_id" cql:"sender_id"`
	RecipientId    gocql.UUID `json:"recipient_id" cql:"recipient_id"`
}

func (c *Conversation) ListMessages(messageId gocql.UUID) ([]ChatMessage, error) {
	var query string
	if len(messageId) == 0 {
		query = fmt.Sprintf(`SELECT *  FROM  messages WHERE conversation_id = %s LIMIT 50`, c.ConversationId)
	} else {
		query = fmt.Sprintf(`SELECT *  FROM  messages WHERE conversation_id = %s and message_id >= %s LIMIT 50`, c.ConversationId, messageId)
	}

	iter := storage.Cassandra.Session.Query(query).Iter()
	scanner := iter.Scanner()
	chatmessages := make([]ChatMessage, iter.NumRows())
	for scanner.Next() {
		var message ChatMessage
		err := scanner.Scan(&message.ConversationId, &message.MessageId, &message.SenderId, &message.RecipientId, &message.Message, &message.Timestamp)
		if err != nil {
			return nil, err
		}
		chatmessages = append(chatmessages, message)
	}
	return chatmessages, nil
}

type ChatMessage struct {
	Conversation
	MessageId gocql.UUID `json:"messaage_id" cql:"message_id"`
	Message   string     `json:"message" cql:"message"`
	Timestamp time.Time  `json:"timestamp" cql:"time_stamp"`
}

func (s *ChatMessage) ModelType() string {
	return "ChatMessage"
}

func (s *ChatMessage) SaveMessage() error {
	query := `INSERT INTO messages(conversation_id,message_id,sender_id,recipient_id,message, time_stamp) 
	VALUES(?, ?,?,?,?,? )`
	err := storage.Cassandra.Session.Query(query, s.ConversationId, s.MessageId, s.SenderId, s.RecipientId, s.Message, s.Timestamp).Exec()

	return err
}

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusUpgradeRequired)
		log.Printf("Failed to upgrade connection to WebSocket: %v", err)
		return
	}

	// Extract user id from the session and register the Conn
	id, _ := gocql.ParseUUID("710cad4e-2ce3-4d81-8e06-59cf5f7d793d")

	user := accounts.User{UserId: id}
	// full user details fetched from DB
	fuser, err := user.GetUser()
	if err != nil {
		return
	}
	client := NewClient(conn, *fuser)

	hub.register <- client

	go client.write(hub)
	client.read(hub)
}
