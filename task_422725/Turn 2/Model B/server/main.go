package main  
import (  
    "bufio"
    "encoding/json"
    "fmt"
    "log"
    "net"
    "strconv"
    "time"
)

type Request struct {    
    Method string `json:"method"`
    Path   string `json:"path"`
}

type Response struct {
    Status int    `json:"status"`
    Body    string `json:"body"`
}

const (
    maxMessageSize = 1024 * 1024 // 1MB
    readTimeout    = time.Second * 5
)

func serveTCP(listenAddr string) {
    log.Printf("Listening on %s", listenAddr)
    ln, err := net.Listen("tcp", listenAddr)
    if err != nil {
        log.Fatalf("Error listening: %v", err)
    }
    defer ln.Close()
    for {
        conn, err := ln.Accept()
        if err != nil {
            log.Printf("Error accepting: %v", err)
            continue
        }
        go handleConn(conn)
    }
}

func handleConn(conn net.Conn) {
    defer conn.Close()
    conn.SetReadDeadline(time.Now().Add(readTimeout))
    for {
        if err := handleMessage(conn); err != nil {
            log.Printf("Error handling message: %v", err)
            break
        }
    }
}

func handleMessage(conn net.Conn) error {
    // Read message length
    lengthBytes := make([]byte, 4)
    _, err := conn.Read(lengthBytes)
    if err != nil {
        return err
    }
    msgLength, err := strconv.Atoi(string(lengthBytes))
    if err != nil {
        return err
    }
    if msgLength > maxMessageSize {
        return fmt.Errorf("message too large: %d bytes", msgLength)
    }
    // Read message body
    msgBuf := make([]byte, msgLength)
    _, err = conn.Read(msgBuf)
    if err != nil {
        return err
    }
    // Parse the JSON request
    var req Request
    if err := json.Unmarshal(msgBuf, &req); err != nil {
        log.Printf("Error parsing request: %v", err)
        return sendResponse(conn, 500, "Invalid request JSON")
    }
    log.Printf("Received request: %+v", req)
    resp := makeResponse(req)
    return sendResponse(conn, resp.Status, resp.Body)
}

func makeResponse(req Request) Response {
    return Response{
        Status: 200,
        Body:   fmt.Sprintf("Hello from path: %s!", req.Path),
    }
}

func sendResponse(conn net.Conn, status int, body string) error {
    resp := Response{
        Status: status,
        Body:   body,
    }
    resBytes, err := json.Marshal(resp)
    if err != nil {
        return err
    }
    // Send response length
    lengthBytes := []byte(strconv.Itoa(len(resBytes)))
    if _, err := conn.Write(lengthBytes); err != nil {
        return err
    }
    // Send response body
    if _, err := conn.Write(resBytes); err != nil {
        return err
    }
    return nil
}

func main() {
    serveTCP(":8080")
}