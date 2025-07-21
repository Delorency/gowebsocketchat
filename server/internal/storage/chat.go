package storage

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn *websocket.Conn
	Name string
}

type chat struct {
	ID int

	clients map[string]*websocket.Conn

	broadcast chan []byte

	registered chan *Client

	unregistered chan *Client
}

func (c *chat) BroadCast(message []byte) {
	for _, conn := range c.clients {
		conn.WriteMessage(websocket.TextMessage, message)
	}
}

func (c *chat) SendToChat(message []byte) {
	c.broadcast <- message
}

func (c *chat) ConnectClient(conn *websocket.Conn, client string) {
	c.registered <- &Client{conn, client}
}
func (c *chat) DisconnectClient(conn *websocket.Conn, client string) {
	c.unregistered <- &Client{conn, client}
}

func (c *chat) IsClientConnected(client string) bool {
	_, ok := c.clients[client]

	return ok
}

func (c *chat) ChatProcessing() {
	for {
		select {
		case client := <-c.registered:
			c.clients[client.Name] = client.Conn
		case client := <-c.unregistered:
			delete(c.clients, client.Name)
			client.Conn.Close()
		case message := <-c.broadcast:
			go c.BroadCast(message)
		}
	}
}
