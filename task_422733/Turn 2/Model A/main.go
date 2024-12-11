package main

import (
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"
)

type trafficConfig struct {
	Server        string
	NumMessages   int
	MessageSize   int
	DelayInterval time.Duration
	TrafficRate   int // messages per second
}

func generateTCPTraffic(config trafficConfig, wg *sync.WaitGroup) {
	defer wg.Done()

	conn, err := net.Dial("tcp", config.Server)
	if err != nil {
		fmt.Println("Error dialing:", err)
		return
	}
	defer conn.Close()

	message := make([]byte, config.MessageSize)

	startTime := time.Now()
	sentMessages := 0

	for sentMessages < config.NumMessages {
		for j := 0; j < config.MessageSize; j++ {
			message[j] = byte(rand.Intn(128))
		}

		_, err := conn.Write(message)
		if err != nil {
			fmt.Println("Error writing:", err)
			return
		}

		sentMessages++

		// Control traffic rate
		if config.TrafficRate > 0 {
			nextSend := time.Now().Add(time.Duration(1000/config.TrafficRate) * time.Millisecond)
			for time.Now() < nextSend {
				time.Sleep(time.Millisecond)
			}
		}

		// Add random delay between messages
		time.Sleep(config.DelayInterval)
	}

	duration := time.Since(startTime)
	fmt.Printf("Sent %d messages in %.2f seconds at %.2f msg/s\n", config.NumMessages, duration.Seconds(), float64(config.NumMessages)/duration.Seconds())
}

func measureServerPerformance(config trafficConfig) {
	wg := &sync.WaitGroup{}
	numGoroutines := config.NumMessages / 1000 // Adjust this based on your system resources

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go generateTCPTraffic(config, wg)
	}

	wg.Wait()
}

func main() {
	rand.Seed(time.Now().UnixNano())

	config1 := trafficConfig{
		Server:        "127.0.0.1:8080",
		NumMessages:   10000,
		MessageSize:   128,
		DelayInterval: 0 * time.Millisecond,
		TrafficRate:   1000,
	}

	fmt.Println("Running test with config:", config1)
	measureServerPerformance(config1)

	config2 := trafficConfig{
		Server:        "127.0.0.1:8080",
		NumMessages:   5000,
		MessageSize:   512,
		DelayInterval: 5 * time.Millisecond,
		TrafficRate:   500,
	}

	fmt.Println("\nRunning test with config:", config2)
	measureServerPerformance(config2)
}