package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"server/internal"
	"server/internal/storage"
	"syscall"
	"time"
)

var host = flag.String("h", "localhost", "IP address")
var port = flag.Int("p", 8080, "Port")

func main() {
	flag.Parse()

	storage := storage.NewStorage()

	server := internal.NewHTTPServer(storage, fmt.Sprintf("%s:%d", *host, *port))

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	go func() {
		log.Printf("Server work on %s:%d\n", *host, *port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln(err.Error())
		}
	}()

	<-ctx.Done()
	fmt.Println("Server is stopping...")
	shtctx, shtcancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shtcancel()

	if err := server.Shutdown(shtctx); err != nil {
		log.Printf("Graceful shutdown error: %v\n", err)

		if err := server.Close(); err != nil {
			log.Fatalf("Forced termination error: %v\n", err)
		}
	}
}
