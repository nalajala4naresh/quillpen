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

func (c *UserConversation) SaveConversation() (*gocql.UUID, error) {
	// add conversation to one user
	// Construct the CQL query dynamically with the bind variables
	// Lookup conversation before saving it
	query := "SELECT conversation_id FROM conversations_by_participants WHERE participants IN ? ;"
	lookupKeys := []string{(c.SenderId.String() + ":" + c.UserId.String()), (c.UserId.String() + ":" + c.SenderId.String())}
	boundQuery := storage.Cassandra.Session.Query(query, lookupKeys)
	resultIter := boundQuery.Iter()

	var conversation_Id gocql.UUID

	if resultIter.Scan(&conversation_Id) {

		return &conversation_Id, nil

	}

	// insert into conversation_by_participants to avoid duplicates
	cbpbatch := storage.Cassandra.Session.NewBatch(gocql.LoggedBatch)
	cbpbatch.Query("INSERT INTO conversations_by_participants(conversation_id,participants) VALUES (?,?)", c.ConversationId, lookupKeys[0])
	cbpbatch.Query("INSERT INTO conversations_by_participants(conversation_id,participants) VALUES (?,?)", c.ConversationId, lookupKeys[1])

	if err := storage.Cassandra.Session.ExecuteBatch(cbpbatch); err != nil {
		fmt.Printf("insert operation in conversations_by_participants failed ")

		return nil, err

	}

	// add the same conversationId into conversations table
	usersbatch := storage.Cassandra.Session.NewBatch(gocql.LoggedBatch)
	
	usersbatch.Query(`INSERT  INTO  conversations(conversation_id,friend_id,friend_name,user_id) VALUES(?,?,?,?)`, c.ConversationId,c.SenderId, c.SenderName,c.UserId)
	usersbatch.Query(` INSERT  INTO  conversations(conversation_id,friend_id,friend_name,user_id) VALUES(?,?,?,?)`, c.ConversationId,c.UserId, c.UserName,c.SenderId)
    
	if err := storage.Cassandra.Session.ExecuteBatch(usersbatch); err != nil {
        fmt.Printf("insert operation in conversations failed ")
		return nil, err

	}

	return &(c.ConversationId), nil
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
		query = fmt.Sprintf(`SELECT conversation_id,message_id,message, sender_id FROM  messages WHERE conversation_id = %s LIMIT 50`, c.ConversationId)
	} else {
		query = fmt.Sprintf(`SELECT conversation_id,message_id,message, sender_id  FROM  messages WHERE conversation_id = %s and message_id <= %s LIMIT 50`, c.ConversationId, messageId)
	}

	iter := storage.Cassandra.Session.Query(query).Iter()
	defer iter.Close()
	scanner := iter.Scanner()
	chatmessages := make([]ConversationMessage, iter.NumRows())
	for scanner.Next() {
		var message ConversationMessage

		err := scanner.Scan(&message.ConversationId, &message.MessageId, &message.Message, &message.SenderId)
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
	VALUES(?, ?,?,? ) USING TTL 86400;`
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

	conv_id, err := conversation.SaveConversation()
	if err != nil {
		log.Printf("Unable to save conversation due to error %s", err)

		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	newConversation := struct {
		ConversationId gocql.UUID `json:"conversation_id"`
	}{ConversationId: *conv_id}

	bcon, err := json.Marshal(newConversation)
	if err != nil {

		log.Printf("Unable marshal saved conversation due to error %s", err)

		w.WriteHeader(http.StatusInternalServerError)
		return

	}
	w.WriteHeader(http.StatusOK)
	w.Write(bcon)

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
