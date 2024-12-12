package main
import (
	"fmt"
	"log"
	"github.com/gorilla/websocket"
)
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer ws.Close()

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Println("unexpected close:", err)
			}
			break
		}
		fmt.Printf("recv: %s\n", message)
	}
}

func main() {
	http.HandleFunc("/ws", serveWs)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("ListenAndServe: %s\n", err)
	}
}