package main
import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http2"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Error listening:", err)
	}
	defer listener.Close()

	s := &http.Server{}
	err = http2.ConfigureServer(s, nil)
	if err != nil {
		log.Fatal("Error configuring HTTP/2:", err)
	}

	go func() {
		if err := s.Serve(listener); err != nil {
			if err == http.ErrServerClosed {
				return
			}
			log.Fatal("Error serving:", err)
		}
	}()

	fmt.Println("Server is running on :8080")
	select {} // Block forever
}