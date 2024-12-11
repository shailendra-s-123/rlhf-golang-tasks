
package main
import (
    "bufio"
    "fmt"
    "log"
    "net"
)
func main() {
    lis, err := net.Listen("tcp", ":1234")
    if err != nil {
        log.Fatal("error listening:", err)
    }
    defer lis.Close()
    for {
        conn, err := lis.Accept()
        if err != nil {
            log.Println("error accepting:", err)
            continue
        }
        defer conn.Close()
        handleConnection(conn)
    }
}
func handleConnection(conn net.Conn) {
    reader := bufio.NewReader(conn)
    writer := bufio.NewWriter(conn)
    for {
        req, err := readRequest(reader)
        if err != nil {
            log.Println("error reading request:", err)
            return
        }
        resp := processRequest(req)
        err = writeResponse(writer, resp)
        if err != nil {
            log.Println("error writing response:", err)
            return
        }
    }
}
