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

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Error upgrading to websocket:", err)
			return
		}
		defer ws.Close()

		for {
			_, message, err := ws.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Println("Error reading websocket:", err)
				}
				return
			}

			err = ws.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Println("Error writing websocket:", err)
				return
			}
		}
	})

	fmt.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}