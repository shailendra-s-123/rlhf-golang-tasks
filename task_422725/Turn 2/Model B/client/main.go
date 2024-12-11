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
    connectTimeout = time.Second * 5
    readTimeout    = time.Second * 5
)

func dialTCP(addr string, path string) Response {
    conn, err := dialWithTimeout(addr)
    if err != nil {
        log.Fatalf("Error dialing: %v", err)
    }
    defer conn.Close()
    req := Request{
        Method: "GET",
        Path:   path,
    }
    reqBytes, err := json.Marshal(req)
    if err != nil {
        log.Printf("Error marshalling request: %v", err)
        return Response{}
    }
    // Send request length
    lengthBytes := []byte(strconv.Itoa(len(reqBytes)))
    if _, err := conn.Write(lengthBytes); err != nil {
        log.Printf("Error writing request length: %v", err)
        return Response{}
    }
    // Send request body
    if _, err := conn.Write(reqBytes); err != nil {
        log.Printf("Error writing request body: %v", err)
        return Response{}
    }
    return readResponse(conn)
}

func dialWithTimeout(addr string) (net.Conn, error) {
    dialer := &net.Dialer{Timeout: connectTimeout}
    conn, err := dialer.Dial("tcp", addr)
    if err != nil {
        return nil, err
    }
    conn.SetReadDeadline(time.Now().Add(readTimeout))
    return conn, nil
}

func readResponse(conn net.Conn) Response {
    // Read response length
    lengthBytes := make([]byte, 4)
    _, err := conn.Read(lengthBytes)
    if err != nil {
        log.Printf("Error reading response length: %v", err)
        return Response{}
    }
    msgLength, err := strconv.Atoi(string(lengthBytes))
    if err != nil {
        log.Printf("Error parsing response length: %v", err)
        return Response{}
    }
    if msgLength > maxMessageSize {
        log.Printf("Response too large: %d bytes", msgLength)
        return Response{}
    }
    // Read response body
    msgBuf := make([]byte, msgLength)
    _, err = conn.Read(msgBuf)
    if err != nil {
        log.Printf("Error reading response body: %v", err)
        return Response{}
    }
    // Parse the JSON response
    var resp Response
    if err := json.Unmarshal(msgBuf, &resp); err != nil {
        log.Printf("Error parsing response: %v", err)
    }
    return resp
}

func main() {
    resp := dialTCP("localhost:8080", "/hello")
    log.Printf("Response: %+v", resp)
}
