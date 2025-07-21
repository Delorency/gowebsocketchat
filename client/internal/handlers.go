package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func CreateChat(host string, port int) {
	url := fmt.Sprintf("http://%s:%d/chats", host, port)
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		log.Fatalf("Request sending error")
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		log.Fatalf("%d | %s", resp.StatusCode, string(body))
	}

	var b struct {
		Id int `json:"id"`
	}

	json.Unmarshal(body, &b)
	log.Printf("Created new chat room: %d", b.Id)
}

func GetChatsList(host string, port int) {
	url := fmt.Sprintf("http://%s:%d/chats", host, port)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Request sending error")
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		log.Fatalf("%d | %s", resp.StatusCode, string(body))
	}
	var b struct {
		Ids []int `json:"list"`
	}

	json.Unmarshal(body, &b)

	log.Printf("List chats has been received: (%d)\n", len(b.Ids))
	for _, v := range b.Ids {
		fmt.Printf("%d ", v)
	}
	fmt.Println()
}

func ConnectToChat(host string, port int, chatID, name string) {
	url := fmt.Sprintf("ws://%s:%d/connect", host, port)

	header := http.Header{}
	header.Set("Chat-ID", chatID)
	header.Set("Client-Name", name)

	conn, _, err := websocket.DefaultDialer.Dial(url, header)
	if err != nil {
		log.Printf("Websocket connection error: %s\n", err.Error())
	}
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Disconnect")
			break
		}
		ClearLastNRows(1)
		fmt.Println(message)
	}

}
