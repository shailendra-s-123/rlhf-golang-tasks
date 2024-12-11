package main  
import (    "bufio"
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

func dialTCP(addr string, path string) Response {
    conn, err := net.Dial("tcp", addr)
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
    fmt.Fprintln(conn, string(reqBytes))
    reader := bufio.NewReader(conn)
    data, err := reader.ReadBytes('\n')
    if err != nil {
        log.Printf("Error reading response: %v", err)
        return Response{}
    }
    // Parse the JSON response
    var resp Response
    if err := json.Unmarshal(data, &resp); err != nil {
        log.Printf("Error parsing response: %v", err)
    }
    return resp
}

func main() {
    resp := dialTCP("localhost:8080", "/hello")
    log.Printf("Response: %+v", resp)
}  