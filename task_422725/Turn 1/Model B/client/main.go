package main
import (
    "bufio"
    "fmt"
    "net"
)

func main() {
    conn, err := net.Dial("tcp", "localhost:1234")
    if err != nil {
        log.Fatal("error dialing:", err)
    }
    defer conn.Close()
    writer := bufio.NewWriter(conn)
    request := createRequest("GET", "/resource", nil)
    err = writeRequest(writer, request)
    if err != nil {
        log.Fatal("error writing request:", err)
    }
    writer.Flush()
    reader := bufio.NewReader(conn)
    response, err := readResponse(reader)
    if err != nil {
        log.Fatal("error reading response:", err)
    }
    fmt.Println(response.Status)
    fmt.Println(response.Body)
}