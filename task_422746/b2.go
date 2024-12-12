package main

import (
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

const (
	maxConnections  = 100
	connectionTTL   = 5 * time.Second
	checkIdlePeriod = 1 * time.Second
)

// ConnPool represents the connection pool
type ConnPool struct {
	mu           sync.Mutex
	activeConnections int32
	connections    map[net.Conn]struct{}
	idleConnections  []net.Conn
	conditions     *sync.Cond
}

// NewConnPool creates a new connection pool
func NewConnPool() *ConnPool {
	cp := &ConnPool{
		connections:    make(map[net.Conn]struct{}),
		idleConnections: []net.Conn{},
		conditions:     sync.NewCond(&cp.mu),
	}
	go cp.cleanupConnections()
	return cp
}

func (cp *ConnPool) Close() {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	for conn := range cp.connections {
		conn.Close()
	}
	for _, conn := range cp.idleConnections {
		conn.Close()
	}
	cp.conditions.Signal()
}

// Dial establishes a new connection or returns an existing one from the pool
func (cp *ConnPool) Dial() (net.Conn, error) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	for {
		if len(cp.idleConnections) > 0 {
			conn := cp.idleConnections[len(cp.idleConnections)-1]
			cp.idleConnections = cp.idleConnections[:len(cp.idleConnections)-1]
			return conn, nil
		}

		if atomic.LoadInt32(&cp.activeConnections) >= maxConnections {
			cp.conditions.Wait()
			continue
		}

		log.Printf("Establishing a new connection...")
		conn, err := net.Dial("tcp", "example.com:80") // Replace with your target address
		if err != nil {
			return nil, err
		}

		atomic.AddInt32(&cp.activeConnections, 1)
		cp.connections[conn] = struct{}{}
		return conn, nil
	}
}

// Return returns a connection to the pool
func (cp *ConnPool) Return(conn net.Conn) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	delete(cp.connections, conn)
	atomic.AddInt32(&cp.activeConnections, -1)

	cp.idleConnections = append(cp.idleConnections, conn)
	cp.conditions.Signal()
}

func (cp *ConnPool) cleanupConnections() {
	for {
		time.Sleep(checkIdlePeriod)
		cp.mu.Lock()
		defer cp.mu.Unlock()

		for _, conn := range cp.idleConnections {
			select {
			case <-time.After(connectionTTL):
				conn.Close()
			default:
				goto next
			}
		}
	next:
		continue
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