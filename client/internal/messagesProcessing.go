package internal

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

func readMessageProcessor(conn *websocket.Conn, cancel context.CancelFunc) {
	defer cancel()
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Disconnect")
			break
		}
		fmt.Println(string(message))
	}
}
func sendMessageProcessor(conn *websocket.Conn, cancel context.CancelFunc) {
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
}
