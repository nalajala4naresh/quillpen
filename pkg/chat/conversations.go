package chat

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"github.com/quillpen/pkg/storage"
)

type Conversation struct {
	FriendId       gocql.UUID `json:"friend_id"`
	FriendName     string     `json:"friend_name"`
	ConversationId gocql.UUID `json:"conversation_id"`
	FriendPublicKey      string     `json:"friend_publickey"`
}

func ListConversations(userId gocql.UUID) ([]Conversation, error) {

	query := "SELECT conversation_id,friend_id,friend_name,friend_publickey FROM  conversations WHERE user_id = ? ;"
	iter := storage.Cassandra.Session.Query(query, userId).Iter()
	scanner := iter.Scanner()
	var results []Conversation
	for scanner.Next() {

		var conversation Conversation
		err := scanner.Scan(&conversation.ConversationId, &conversation.FriendId, &conversation.FriendName,&conversation.FriendPublicKey)
		if err != nil {
			return nil, err
		}
		results = append(results, conversation)

	}

	return results, nil

}

// chat page List all the conversations a user has
func ListConversationsHandler(w http.ResponseWriter, r *http.Request) {

	vals := mux.Vars(r)
	userId := vals["userId"]
	id, err := gocql.ParseUUID(userId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	conversations, err := ListConversations(id)

	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)

	}

	data, err := json.Marshal(conversations)
	log.Printf("length of data is %d", len(data))
	if err != nil {
		log.Printf("Json Marshalling failed for conversations")
		w.WriteHeader(http.StatusInternalServerError)

	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)

}

// Delete a conversation in the above list from user

func DeleteConversation(convId gocql.UUID) error {

	// Delete conversation
	query := storage.Cassandra.Session.Query("DELETE FROM messages WHERE  conversation_id = ?", convId)
	return query.Exec()

}
func DeleteConversationHandler(w http.ResponseWriter, r *http.Request) {

	vals := mux.Vars(r)
	convId := vals["conversationId"]
	cId, err := gocql.ParseUUID(convId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = DeleteConversation(cId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	w.WriteHeader(http.StatusOK)

}
