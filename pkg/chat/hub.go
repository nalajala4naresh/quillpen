package chat

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/gocql/gocql"
	"github.com/gorilla/websocket"

	"github.com/quillpen/pkg/accounts"
)

func NewConversationChanel(conn *websocket.Conn, user accounts.User) *ConversationChanel {
	return &ConversationChanel{conn: conn, send: make(chan ConversationMessage)}
}

type ConversationChanel struct {
	conn           *websocket.Conn
	userId         gocql.UUID
	conversationId gocql.UUID
	send           chan ConversationMessage
}

func (c *ConversationChanel) read(hub *Hub) {
	defer func() {
		c.conn.Close()
		hub.unregister <- c
	}()

	for {

		var mess ConversationMessage

		_, rawbytes, err := c.conn.ReadMessage()
		if err != nil {
			fmt.Printf("%s", err)
			break
		}

		err = json.Unmarshal(rawbytes, &mess)

		if err != nil {
			// Invalid chat message format
			log.Printf("JSON error is %s", err)
			break
		}
		fmt.Print(string(rawbytes))
		// extract the conversation from message and write it to the cassaandra
		hub.broadcast <- mess
	}
}

func (c *ConversationChanel) write(hub *Hub) {
	defer c.conn.Close()

	for message := range c.send {
		bmess, err := json.Marshal(message)
		if err != nil {
			// json Marshalling error
			log.Printf("%s", err)
			break
		}

		err = c.conn.WriteMessage(websocket.TextMessage, bmess)

		if err != nil {
			// data loss due to channels
			break
		}

	}
}

type Hub struct {
	conversations map[gocql.UUID]map[gocql.UUID]*ConversationChanel
	register      chan *ConversationChanel
	unregister    chan *ConversationChanel
	broadcast     chan ConversationMessage
	mu            sync.Mutex
}

func (h *Hub) run() {
	for {
		select {
		case connectionChannel := <-h.register:
			// add the new connection to the channel
			h.mu.Lock()
			participants := h.conversations[connectionChannel.conversationId]

			// add a participant
			if participants == nil {
				participants = make(map[gocql.UUID]*ConversationChanel)

			}
			participants[connectionChannel.userId] = connectionChannel

			h.conversations[connectionChannel.conversationId] = participants
			h.mu.Unlock()

			conversation := ConversationMessage{ConversationId: connectionChannel.conversationId } 
            messages , _ := conversation.ListMessages(nil)
			for _, message := range messages {
				h.broadcast <- message
			}
		case connectionChannel := <-h.unregister:
			h.mu.Lock()
			participants := h.conversations[connectionChannel.conversationId]

			h.conversations[connectionChannel.conversationId] = participants

			h.mu.Unlock()
			if _, ok := participants[connectionChannel.userId]; ok {
				close(connectionChannel.send)
				delete(h.conversations[connectionChannel.conversationId], connectionChannel.userId)

			}
		case message := <-h.broadcast:
			conversationId := message.ConversationId
			// take all the participants
			participants := h.conversations[conversationId]
            
			for _, conn := range participants {

				conn.send <- message

			}

			fmt.Printf("%s, %s, %s,%s",message.ConversationId, message.MessageId, message.SenderId, message.Message)

			// write the message to database
			err := message.SaveMessage()
			if err != nil {
				// data loss so closing websocket
				// cassandra error
				log.Printf("%s", err)
				delete(h.conversations, conversationId)

			}

		}
	}
}

func NewHub() *Hub {
	return &Hub{
		conversations: make(map[gocql.UUID]map[gocql.UUID]*ConversationChanel, 10000),
		register:      make(chan *ConversationChanel, 10000),
		unregister:    make(chan *ConversationChanel, 10000),
		broadcast:     make(chan ConversationMessage, 100000),
	}
}
