package main
import (
	"fmt"
	"math/rand"
	"net"
	"time"
)
// SimulateClient generates random traffic to a server
func SimulateClient(addr string, rate float64, duration time.Duration) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("Error dialing:", err)
		return
	}
	defer conn.Close()

	ticker := time.NewTicker(time.Duration(1000.0 / rate) * time.Millisecond)
	defer ticker.Stop()

	start := time.Now()
	for {
		select {
		case <-ticker.C:
			data := make([]byte, rand.Intn(100)+10) // Generate random payload between 10 and 109 bytes
			n, err := conn.Write(data)
			if err != nil {
				fmt.Println("Error writing:", err)
				return
			}
			fmt.Printf("Sent %d bytes\n", n)

			// Read response (if any)
			resp := make([]byte, 1024)
			n, err = conn.Read(resp)
			if err != nil {
				fmt.Println("Error reading:", err)
				return
			}
			fmt.Printf("Received %d bytes: %s\n", n, string(resp[:n]))

		case <-time.After(duration):
			fmt.Println("Client simulation completed.")
			return
		}
	}
}
func main() {
	rand.Seed(time.Now().UnixNano())
	serverAddr := "127.0.0.1:8080" // Replace this with your server's address
	trafficRate := 100.0 // 100 packets per second
	simulationDuration := 10 * time.Second

	go SimulateClient(serverAddr, trafficRate, simulationDuration)
	time.Sleep(simulationDuration + 1*time.Second)
}