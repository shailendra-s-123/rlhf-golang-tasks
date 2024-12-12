package main

import (
	"log"
	"websocket"
)

func main() {
	ws, err := websocket.Dial("ws://example.com/ws", nil, nil)
	if err != nil {
		log.Fatal("Error dialing:", err)
	}
	defer ws.Close()

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("Connection closed:", err)
			}
			return
		}
		log.Println("Received:", string(message))
	}
}