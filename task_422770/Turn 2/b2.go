package main

import (
	"fmt"
	"log"
	"net/http"
	"websocket"
)

func main() {
	url := "ws://example.com/ws"

	ws, err := websocket.Dial(url, "", nil)
	if err != nil {
		log.Fatal("Connection failed:", err)
	}
	defer ws.Close()

	// Rest of the WebSocket communication logic
}