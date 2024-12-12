
package main

import (
	"net"
	"sync"
)

const (
	maxConnections = 100 // Maximum number of connections allowed
	connectionTimeout = 10 * 1000 // Connection timeout in milliseconds
)

type connectionPool struct {
	mu       sync.Mutex
	connections []*net.TCPConn
	max       int
	timeout   int
}

func newConnectionPool(max int, timeout int) *connectionPool {
	return &connectionPool{
		connections: make([]*net.TCPConn, 0),
		max:         max,
		timeout:     timeout,
	}
}

func (cp *connectionPool) get() (*net.TCPConn, bool) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	for _, conn := range cp.connections {
		// Check if the connection is still active
		if _, err := conn.Write([]byte(" ")); err != nil {
			cp.remove(conn)
			continue
		}
		return conn, true
	}

	return nil, false
}

func (cp *connectionPool) put(conn *net.TCPConn) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	if len(cp.connections) >= cp.max {
		// Close the connection if the pool is full
		conn.Close()
		return
	}

	cp.connections = append(cp.connections, conn)
}

func (cp *connectionPool) remove(conn *net.TCPConn) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	for i, c := range cp.connections {
		if c == conn {
			cp.connections = append(cp.connections[:i], cp.connections[i+1:]...)
			conn.Close()
			return
		}
	}
}

func tcpConnectionMiddleware(handler func(*net.TCPConn)) {
	pool := newConnectionPool(maxConnections, connectionTimeout)

	go func() {
		for {
			conn, err := net.Listen("tcp", ":12345")
			if err != nil {
				panic(err)
			}
			defer conn.Close()

			for {
				newConn, err := conn.Accept()
				if err != nil {
					continue
				}

				go func(c *net.TCPConn) {
					if conn, ok := pool.get(); ok {
						defer pool.put(conn)
						// Handle request using the recycled connection
						handler(conn)
					} else {
						// Handle request using a new connection
						defer c.Close()
						handler(c)
					}
				}(newConn)
			}
		}
	}()
}

func main() {
	tcpConnectionMiddleware(func(conn *net.TCPConn) {
		// Your request handling logic here
		_, err := conn.Write([]byte("Hello, client!\n"))
		if err != nil {
			panic(err)
		}
	})
}
