package main

import (
	"fmt"
	"log"
	"net/http"
	"websocket"
)

func main() {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	url := "ws://example.com/ws"

	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Println("Error connecting to WebSocket:", err)
		return
	}
	defer ws.Close()

	// Rest of the application logic
}