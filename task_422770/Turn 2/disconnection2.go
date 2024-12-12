package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"websocket"
)

func main() {
	url := "ws://example.com/ws"

	ws, err := websocket.Dial(url, "", nil)
	if err != nil {
		log.Fatal("Connection failed:", err)
	}
	defer ws.Close()

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("Disconnected:", err)
			} else {
				log.Println("Error reading message:", err)
			}
			return
		}
		fmt.Println("Received:", string(message))

		time.Sleep(1 * time.Second)

		err = ws.WriteMessage(websocket.TextMessage, []byte("Hello"))
		if err != nil {
			log.Println("Error writing message:", err)
			return
		}
	}
}