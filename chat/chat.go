package chat

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"github.com/quillpen/sessionManager"
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
	MessageId      gocql.UUID `json:"message_id" cql:"message_id"`
	Message        string     `json:"message" cql:"message"`
}

func (c *Conversation) ListMessages(messageId gocql.UUID) ([]Conversation, error) {
	var query string
	if len(messageId) == 0 {
		query = fmt.Sprintf(`SELECT *  FROM  conversations WHERE conversation_id = %s LIMIT 50`, c.ConversationId)
	} else {
		query = fmt.Sprintf(`SELECT *  FROM  conversations WHERE conversation_id = %s and message_id >= %s LIMIT 50`, c.ConversationId, messageId)
	}

	iter := storage.Cassandra.Session.Query(query).Iter()
	scanner := iter.Scanner()
	chatmessages := make([]Conversation, iter.NumRows())
	for scanner.Next() {
		var message Conversation
		err := scanner.Scan(&message.ConversationId, &message.MessageId, &message.SenderId, &message.Message)
		if err != nil {
			return nil, err
		}
		chatmessages = append(chatmessages, message)
	}
	return chatmessages, nil
}

func (s *Conversation) SaveMessage() error {
	query := `INSERT INTO conversations(conversation_id,message_id,sender_id,message) 
	VALUES(?, ?,?,? )`
	err := storage.Cassandra.Session.Query(query, s.ConversationId, s.MessageId, s.SenderId, s.Message).Exec()

	return err
}

func ConversationsHandler(w http.ResponseWriter, R *http.Request) {

}

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	var userId gocql.UUID

	session, _ := sessionManager.Store.Get(r, sessionManager.SessionName)
	if session.IsNew {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
	} else {
		suserId := session.Values[sessionManager.SessionUserId].(string)
		userId, _ = gocql.ParseUUID(suserId)

	}

	vals := mux.Vars(r)
	conversationId := vals["id"]
	convUuid, err := gocql.ParseUUID(conversationId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return

	}

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

	conchanel := ConversationChanel{conn: conn, userId: userId, conversationId: convUuid, send: make(chan Conversation)}
	// full user details fetched from DB

	hub.register <- &conchanel

	go conchanel.write(hub)
	conchanel.read(hub)
}
