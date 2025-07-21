package storage

import (
	"sync"

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

	mu sync.RWMutex
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

func (c *chat) ChatProcessing(storage *Storage) {
	for {
		select {
		case client := <-c.registered:
			c.mu.Lock()
			c.clients[client.Name] = client.Conn
			c.mu.Unlock()
		case client := <-c.unregistered:
			c.mu.Lock()
			delete(c.clients, client.Name)
			c.mu.Unlock()

			if len(c.clients) == 0 {
				storage.mu.Lock()
				delete(storage.chats, c.ID)
				storage.index--
				storage.mu.Unlock()
			}

			client.Conn.Close()
		case message := <-c.broadcast:
			go c.BroadCast(message)
		}
	}
}
