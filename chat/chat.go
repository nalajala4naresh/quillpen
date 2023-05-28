package chat

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gobwas/ws"
	"github.com/gocql/gocql"

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

func (c *Conversation) ListMessages(timestamp *time.Time) ([]ChatMessage, error) {
	var query string
	if timestamp == nil {
		query = fmt.Sprintf(`SELECT *  FROM  TABLE messages WHERE conversation_id = %s LIMIT 50`, c.ConversationId)
	} else {
		query = fmt.Sprintf(`SELECT *  FROM  TABLE messages WHERE conversation_id = %s and time_stamp >= %s LIMIT 50`, c.ConversationId, timestamp)
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
	if len(s.ConversationId) == 0 {
		s.ConversationId = gocql.MustRandomUUID()
	}
	s.MessageId = gocql.TimeUUID()

	query := `INSERT INTO TABLE messages(conversation_id,message_id,sender_id,recipient_id,message, time_stamp) 
	VALUES(?, ?,?,?,?,? )`
	err := storage.Cassandra.Session.Query(query, s.ConversationId, s.MessageId, s.SenderId, s.RecipientId, s.Message, s.Timestamp).Exec()

	return err
}

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		// Unable to upgrade the request
		w.WriteHeader(http.StatusUpgradeRequired)
		return

	}
	// Extract user id from the session and register the Conn
	var id gocql.UUID

	client := NewClient(conn, id)
	hub.register <- client

	go client.write(hub)
	client.read(hub)
}
