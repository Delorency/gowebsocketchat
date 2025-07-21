package main

import (
	"client/internal"
	"flag"
	"log"
	"strconv"
)

var host = flag.String("h", "localhost", "Server IP address")
var port = flag.Int("p", 8080, "Server port")
var name = flag.String("n", "", "Your name")
var chat = flag.Int("chat", 0, "Chat id")
var list = flag.Bool("list", false, "Show chats list")

func main() {
	flag.Parse()
	if *name == "" {
		log.Fatalln("Run command with `-n` to register")
	}
	if *list == true {
		internal.GetChatsList(*host, *port)

	} else if *chat == 0 {
		internal.CreateChat(*host, *port)
	} else if *chat > 0 {
		n := strconv.Itoa(*chat)
		internal.ConnectToChat(*host, *port, n, *name)
	}

}
