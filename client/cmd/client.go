package main

import (
	"client/internal"
	"flag"
	"log"
)

var host = flag.String("h", "localhost", "Server IP address")
var port = flag.Int("p", 8080, "Server port")
var name = flag.String("n", "", "Your name")
var chat = flag.Int("chat", 0, "Chat id")

func main() {
	flag.Parse()
	if *name == "" {
		log.Fatalln("Run command with `-n` to register")
	}

	if *chat == 0 {
		internal.CreateChat(*host, *port)
	}

}
