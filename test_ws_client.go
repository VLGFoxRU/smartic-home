package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Usage: go run test_ws_client.go <url> <token>")
	}
	url := os.Args[1]
	token := os.Args[2]

	header := http.Header{}
	header.Set("Authorization", "Bearer "+token)
	// Но мы передаём токен в query, поэтому url уже содержит ?token=
	conn, _, err := websocket.DefaultDialer.Dial(url, header)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		log.Printf("recv: %s", message)
	}
}