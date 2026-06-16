package ws

import (
	"log"

	"github.com/gorilla/websocket"
	userType "github.com/umesshk/termi-chatt/internal/user"
)

type Client struct {
	Conn *websocket.Conn
	Send chan userType.ServerResponse
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		Conn: conn,
		Send: make(chan userType.ServerResponse, 256),
	}
}

func (c *Client) WritePump() {
	defer c.Conn.Close()
	for resp := range c.Send {
		if err := c.Conn.WriteJSON(resp); err != nil {
			log.Println("write error:", err)
			return
		}
	}
}

func (c *Client) Enqueue(resp userType.ServerResponse) {
	select {
	case c.Send <- resp:
	default:
		log.Println("client send buffer full, dropping message")
	}
}
