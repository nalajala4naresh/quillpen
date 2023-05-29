package chat

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gocql/gocql"
	"github.com/gorilla/websocket"

	"github.com/quillpen/accounts"
)

func NewClient(conn *websocket.Conn, id gocql.UUID) *Client {
	return &Client{conn: conn, send: make(chan ChatMessage), userId: id}
}

type Client struct {
	conn   *websocket.Conn
	send   chan ChatMessage
	userId gocql.UUID
}

func (c *Client) read(hub *Hub) {
	defer func() {
		c.conn.Close()
		hub.unregister <- c
	}()

	for {

		var mess ChatMessage

		_, rawbytes, err := c.conn.ReadMessage()
		if err != nil {
			fmt.Printf("%s", err)
			break
		}

		err = json.Unmarshal(rawbytes, &mess)
		mess.Timestamp = time.Now()

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

func (c *Client) write(hub *Hub) {
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
	clients    map[gocql.UUID]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan ChatMessage
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.userId] = client
			// get user conversations list and read from last read message_id and feed it into send on each client.
			user := accounts.User{UserId: client.userId}
			cuser, err := user.GetUser()
			if err != nil {
				log.Fatalf("Unable to Get User from DB with error %s", err)
				break
			}

			for conv_id, message_id := range cuser.Conversations {
				conversation := Conversation{ConversationId: conv_id}
				messages, err := conversation.ListMessages(message_id)
				if err != nil {
					break
				}
				for _, mess := range messages {
					h.broadcast <- mess
				}

			}

		case client := <-h.unregister:
			if _, ok := h.clients[client.userId]; ok {
				close(client.send)
				delete(h.clients, client.userId)

			}
		case message := <-h.broadcast:
			client, ok := h.clients[message.RecipientId]
			if ok {
				client.send <- message
			}

			// write the message to database
			err := message.SaveMessage()
			if err != nil {
				// data loss so closing websocket
				// cassandra error
				log.Printf("%s", err)
				close(client.send)
				delete(h.clients, client.userId)

			}

		}
	}
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[gocql.UUID]*Client, 10000),
		register:   make(chan *Client, 10000),
		unregister: make(chan *Client, 10000),
		broadcast:  make(chan ChatMessage, 256),
	}
}
