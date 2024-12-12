
package main

import (
	"log"
	"net"
	"sync"
	"time"
)

const (
	maxConnections = 10
	connectionTTL  = 5 * time.Second
)

// ConnPool represents the connection pool
type ConnPool struct {
	mu           sync.Mutex
	connections chan net.Conn
	conditions  *sync.Cond
	closed       bool
}

// NewConnPool creates a new connection pool
func NewConnPool() *ConnPool {
	cp := &ConnPool{
		connections: make(chan net.Conn, maxConnections),
		conditions:  sync.NewCond(&cp.mu),
	}
	go cp.cleanupConnections()
	return cp
}

func (cp *ConnPool) Close() {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	cp.closed = true
	cp.conditions.Signal()
	close(cp.connections)
}

// Dial establishes a new connection or returns an existing one from the pool
func (cp *ConnPool) Dial() (net.Conn, error) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	for {
		select {
		case conn := <-cp.connections:
			return conn, nil // Reuse an existing connection if available
		default:
		}

		// Limit the number of active connections
		if len(cp.connections) == maxConnections {
			cp.conditions.Wait()
			continue
		}

		log.Printf("Establishing a new connection...")
		conn, err := net.Dial("tcp", "example.com:80") // Replace with your target address
		if err != nil {
			return nil, err
		}

		return conn, nil
	}
}

// Return returns a connection to the pool
func (cp *ConnPool) Return(conn net.Conn) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	if cp.closed {
		conn.Close()
		return
	}

	go func(conn net.Conn) {
		cp.conditions.L.Lock()
		defer cp.conditions.L.Unlock()

		select {
		case <-time.After(connectionTTL):
			conn.Close()
		case cp.connections <- conn:
		}
		cp.conditions.Signal()
	}(conn)
}

func (cp *ConnPool) cleanupConnections() {
	for {
		cp.mu.Lock()
		for len(cp.connections) == 0 {
			if cp.closed {
				cp.mu.Unlock()
				return
			}
			cp.conditions.Wait()
		}
		cp.mu.Unlock()

		select {
		case <-time.After(1 * time.Second):
		default:
		}
	}
}

func main() {
	cp := NewConnPool()

	go func() {
		for {
			conn, err := cp.Dial()
			if err != nil {
				log.Fatalf("Error dialing: %v", err)
			}

			go func(conn net.Conn) {
				// Simulate work using the connection
				time.Sleep(2 * time.Second)
				_, err := conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\nHello, World!"))
				if err != nil {
					log.Printf("Error writing to connection: %v", err)
				}
				cp.Return(conn)
			}(conn)
		}
	}()

	select {}
}
