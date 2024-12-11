package main  

import (
	"fmt"
	"math/rand"
	"net"
	"time"
)

func generateTCPTraffic(server string, numMessages int) {
	conn, err := net.Dial("tcp", server)
	if err != nil {
		fmt.Println("Error dialing:", err)
		return
	}
	defer conn.Close()

	message := "Hello, Server!"
	for i := 0; i < numMessages; i++ {
		_, err := fmt.Fprintln(conn, message)
		if err != nil {
			fmt.Println("Error writing:", err)
			return
		}
		// Add some randomness to the traffic generation
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	}
	fmt.Println("Sent", numMessages, "messages to", server)
}

func main() {
	server := "127.0.0.1:8080" // Replace this with your server's address
	numMessages := 1000
	generateTCPTraffic(server, numMessages)
}