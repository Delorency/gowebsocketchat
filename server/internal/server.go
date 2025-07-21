package internal

import (
	"net/http"
	"server/internal/storage"
)

func NewHTTPServer(storage storage.IStorage, addr string) *http.Server {
	router := http.NewServeMux()

	handler := NewHandler(storage)

	router.HandleFunc("/chats", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			handler.ListChats(w, r)
		case http.MethodPost:
			handler.CreateNewChat(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	})
	router.HandleFunc("/ws/connect", handler.SendMessage)

	return &http.Server{Addr: addr, Handler: router}
}
