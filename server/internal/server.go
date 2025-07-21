package internal

import (
	"net/http"
	"server/internal/storage"
)

func NewHTTPServer(storage storage.IStorage, addr string) *http.Server {
	router := http.NewServeMux()

	handler := NewHandler(storage)

	router.HandleFunc("/registered", handler.CreateNewChat)
	router.HandleFunc("/ws/send", handler.SendMessage)

	return &http.Server{Addr: addr, Handler: router}
}
