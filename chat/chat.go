package chat

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"github.com/quillpen/storage"
)

// chat to store send and recieve messages
var hub *Hub

func init() {
	//  create keyspace and tables

	hub = NewHub()
	go hub.run()
}

type UserConversation struct {
	ConversationId gocql.UUID `json:"conversation_id"`
	SenderId       gocql.UUID `json:"sender_id"`
	UserId         gocql.UUID `json:"user_id"`
	UserName       string     `json:"user_name"`
	SenderName     string     `json:"sender_name"`
}

func (c *UserConversation) SaveConversation() error {
	// add conversation to one user
	// Construct the CQL query dynamically with the bind variables
	query := "INSERT INTO conversations(conversation_id,participants) VALUES (?,?)"

	if err := storage.Cassandra.Session.Query(query, c.ConversationId, []gocql.UUID{c.SenderId, c.UserId}).Exec(); err != nil {
		return err
	}

	batch := storage.Cassandra.Session.NewBatch(gocql.UnloggedBatch)
	batch.Query(`UPDATE users SET conversations = conversations +  ? WHERE user_id = ?;`, []gocql.UUID{c.ConversationId}, c.UserId)
	batch.Query(` UPDATE users SET conversations = conversations +  ? WHERE user_id = ?;`, []gocql.UUID{c.ConversationId}, c.SenderId)

	if err := storage.Cassandra.Session.ExecuteBatch(batch); err != nil {

		return err

	}

	return nil
}

type ConversationMessage struct {
	ConversationId gocql.UUID `json:"conversation_id" cql:"conversation_id"`
	SenderId       gocql.UUID `json:"sender_id" cql:"sender_id"`
	MessageId      gocql.UUID `json:"message_id" cql:"message_id"`
	Message        string     `json:"message" cql:"message"`
}

func (c *ConversationMessage) ListMessages(messageId *gocql.UUID) ([]ConversationMessage, error) {
	var query string
	if messageId == nil {
		query = fmt.Sprintf(`SELECT conversation_id,message_id,message, sender_id FROM  messages WHERE conversation_id = %s LIMIT 10`, c.ConversationId)
	} else {
		query = fmt.Sprintf(`SELECT conversation_id,message_id,message, sender_id  FROM  messages WHERE conversation_id = %s and message_id <= %s LIMIT 50`, c.ConversationId, messageId)
	}

	iter := storage.Cassandra.Session.Query(query).Iter()
	defer iter.Close()
	scanner := iter.Scanner()
	chatmessages := make([]ConversationMessage, iter.NumRows())
	for scanner.Next() {
		var message ConversationMessage

		err := scanner.Scan(&message.ConversationId, &message.MessageId, &message.Message,&message.SenderId)
		if err != nil {
			fmt.Printf("%s", err)
			return nil, err
		}
		
		chatmessages = append(chatmessages, message)
	}

	return chatmessages, nil
}

func (s *ConversationMessage) SaveMessage() error {
	query := `INSERT INTO messages(conversation_id,message_id,sender_id,message) 
	VALUES(?, ?,?,? )`
	err := storage.Cassandra.Session.Query(query, s.ConversationId, s.MessageId, s.SenderId, s.Message).Exec()

	return err
}

func ConversationsHandler(w http.ResponseWriter, r *http.Request) {

	var conversation UserConversation
	defer r.Body.Close()

	jb, _ := io.ReadAll(r.Body)

	jerr := json.Unmarshal(jb, &conversation)
	if jerr != nil {
		log.Printf("Unable to Unmarshall due to error %s", jerr)
		w.WriteHeader(http.StatusInternalServerError)
	}

	err := conversation.SaveConversation()
	if err != nil {
		log.Printf("Unable to save conversation due to error %s", err)

		w.WriteHeader(http.StatusInternalServerError)
		return

	}
	w.WriteHeader(http.StatusAccepted)

}

func ChatHandler(w http.ResponseWriter, r *http.Request) {

	// session, _ := sessionManager.Store.Get(r, sessionManager.SessionName)
	// if session.IsNew {
	// 	http.Redirect(w, r, "/signin", http.StatusSeeOther)
	// } else {
	// 	suserId := session.Values[sessionManager.SessionUserId].(string)
	// 	userId, _ = gocql.ParseUUID(suserId)

	// }

	vals := mux.Vars(r)
	conversationId, err := gocql.ParseUUID(vals["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return

	}
	userId, err := gocql.ParseUUID(vals["userid"])
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

	conchanel := ConversationChanel{conn: conn, userId: userId, conversationId: conversationId, send: make(chan ConversationMessage)}
	// full user details fetched from DB

	hub.register <- &conchanel

	go conchanel.write(hub)
	conchanel.read(hub)
}
