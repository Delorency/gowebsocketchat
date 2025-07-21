package storage

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type IStorage interface {
	AddChat() (int, error)
	GetChat(id int) *chat
	ListChats() []int
}

type Storage struct {
	chats map[int]*chat
	mu    sync.RWMutex
	index int
	size  int
}

func (s *Storage) AddChat() (int, error) {
	if s.index > s.size {
		return 0, fmt.Errorf("Maximum chat quantity (%d)", s.size)
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

func (s *Storage) ListChats() []int {
	list := make([]int, s.size)
	i := 0
	for index, _ := range s.chats {
		list[i] = index
		i++
	}
	return list[:i]
}

func NewStorage(size int) IStorage {
	return &Storage{chats: make(map[int]*chat), index: 1, size: size}
}
