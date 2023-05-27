package chat

import (
	"net"

	"github.com/gobwas/ws/wsutil"
)

func NewClient(conn net.Conn) *Client {
	return &Client{conn: conn, send: make(chan []byte)}
}

type Client struct {
	conn net.Conn
	send chan []byte
}

func (c *Client) read(hub *Hub) {
	defer func() {
		c.conn.Close()
		hub.unregister <- c
	}()

	for {

		message, _, err := wsutil.ReadClientData(c.conn)
		if err != nil {
			break
		}
		hub.broadcast <- message
	}
}

func (c *Client) write(hub *Hub) {
	defer c.conn.Close()

	for message := range c.send {

		err := wsutil.WriteServerBinary(c.conn, message)
		if err != nil {
			// data loss due to channels
			break
		}

	}
}

type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				close(client.send)
				delete(h.clients, client)

			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool, 10000),
		register:   make(chan *Client, 10000),
		unregister: make(chan *Client, 10000),
		broadcast:  make(chan []byte, 256),
	}
}
