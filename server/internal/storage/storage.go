package storage

import (
	"time"

	"github.com/gorilla/websocket"
)

type IStorage interface {
	AddChat() int
	GetChat(id int) *chat
}

type storage struct {
	chats map[int]*chat
}

func (s *storage) AddChat() int {
	id := time.Now().Second()
	s.chats[id] = &chat{
		ID:           id,
		clients:      make(map[string]*websocket.Conn),
		broadcast:    make(chan []byte),
		registered:   make(chan *Client),
		unregistered: make(chan *Client),
	}

	go s.chats[id].ChatProcessing()

	return id
}

func (s *storage) GetChat(id int) *chat {
	chatptr, ok := s.chats[id]
	if !ok {
		return nil
	}
	return chatptr
}

func NewStorage() IStorage {
	return &storage{chats: make(map[int]*chat)}
}
