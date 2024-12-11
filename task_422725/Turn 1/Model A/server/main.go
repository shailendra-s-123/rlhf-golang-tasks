package main

import (
    "bufio"
    "encoding/json"
    "fmt"
    "log"
    "net"
)
type Request struct {
    Method string `json:"method"`
    Path   string `json:"path"`
}

type Response struct {
    Status int    `json:"status"`
    Body    string `json:"body"`
}

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
        handleConn(conn)
    }
}

func handleConn(conn net.Conn) {
    defer conn.Close()
    reader := bufio.NewReader(conn)
    for {
        data, err := reader.ReadBytes('\n')
        if err != nil {
            if err == io.EOF {
                log.Println("Connection closed")
                break
            }
            log.Printf("Error reading: %v", err)
            continue
        }
        // Parse the JSON request
        var req Request
        if err := json.Unmarshal(data, &req); err != nil {
            log.Printf("Error parsing request: %v", err)
            sendResponse(conn, 500, "Invalid request JSON")
            continue
        }
        log.Printf("Received request: %+v", req)
        resp := makeResponse(req)
        sendResponse(conn, resp.Status, resp.Body)
    }
}

func makeResponse(req Request) Response {
    return Response{
        Status: 200,
        Body:   fmt.Sprintf("Hello from path: %s!", req.Path),
    }
}

func sendResponse(conn net.Conn, status int, body string) {
    resp := Response{
        Status: status,
        Body:   body,
    }
    resBytes, err := json.Marshal(resp)
    if err != nil {
        log.Printf("Error marshalling response: %v", err)
        return
    }
    fmt.Fprintln(conn, string(resBytes))
}
func main() {
    serveTCP(":8080")
}