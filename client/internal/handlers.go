package internal

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

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

	conn, resp, err := websocket.DefaultDialer.Dial(url, header)
	if err != nil {
		body, _ := io.ReadAll(resp.Body)

		log.Fatalf("Websocket connection error: %s\n", body)
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer cancel()
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Disconnect")
				break
			}
			fmt.Println(string(message))
		}
	}()
	go func() {
		defer cancel()
		reader := bufio.NewReader(os.Stdin)
		for {
			input, err := reader.ReadString('\n')
			if err != nil {
				continue
			}
			input = strings.TrimSpace(input)
			ClearLastNRows(1)

			err = conn.WriteMessage(websocket.TextMessage, []byte(input))
			if err != nil {
				log.Println("Send message error")
				break
			}
		}
	}()

	<-ctx.Done()
}
