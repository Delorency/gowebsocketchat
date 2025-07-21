package storage

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type IStorage interface {
	AddChat() (int, error)
	GetChat(id int) *chat
}

type Storage struct {
	chats map[int]*chat
	mu    sync.RWMutex
	index int
}

func (s *Storage) AddChat() (int, error) {
	if s.index > 60 {
		return 0, fmt.Errorf("Maximum chat quantity (60)")
	}
	s.mu.Lock()
	id := s.index
	s.chats[id] = &chat{
		ID:           id,
		clients:      make(map[string]*websocket.Conn),
		broadcast:    make(chan []byte),
		registered:   make(chan *Client),
		unregistered: make(chan *Client),
	}
	s.index++
	s.mu.Unlock()

	go s.chats[id].ChatProcessing(s)

	return id, nil
}

func (s *Storage) GetChat(id int) *chat {
	chatptr, ok := s.chats[id]
	if !ok {
		return nil
	}
	return chatptr
}

func NewStorage() IStorage {
	return &Storage{chats: make(map[int]*chat), index: 1}
}
