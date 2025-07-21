package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func CreateChat(host string, port int) {
	url := fmt.Sprintf("http://%s:%d/registered", host, port)
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
