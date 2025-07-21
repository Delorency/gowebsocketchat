package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"server/internal/storage"
	"strconv"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type IHandler interface {
	CreateNewChat(w http.ResponseWriter, r *http.Request)
	SendMessage(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	storage storage.IStorage
}

func NewHandler(storage storage.IStorage) IHandler {
	return &handler{storage: storage}
}

func (h *handler) CreateNewChat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id, err := h.storage.AddChat()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]int{"id": id})

	log.Printf("Created new chat room with id: %d", id)
}

func (h *handler) SendMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(r.Header.Get("Chat-ID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("'Chat-ID' header does not exist"))
		return
	}
	name := r.Header.Get("Chat-Name")
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("'Chat-Name' header does not exist"))
		return
	}

	chat := h.storage.GetChat(id)
	if chat == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Chat does not exist"))
		return
	}
	connected := chat.IsClientConnected(name)
	if connected {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("The client has already connected"))
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error to connect client")
		return
	}
	go chat.ConnectClient(conn, name)

	for {
		mt, message, err := conn.ReadMessage()
		if err != nil || mt == websocket.CloseMessage {
			break
		}

		b := bytes.NewBufferString(fmt.Sprintf("%s: ", name))
		_, err = b.Write(message)

		if err != nil {
			err1 := conn.WriteMessage(websocket.TextMessage, []byte("Cannot send message"))
			if err1 != nil {
				log.Println("write close:", err1)
				break
			}
		}

		go chat.SendToChat(b.Bytes())
	}

	go chat.DisconnectClient(conn, name)
}
