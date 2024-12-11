
package main
import (
    "encoding/binary"
    "encoding/json"
    "io"
)
func writeRequest(writer io.Writer, req *Request) error {
    data, err := json.Marshal(req)
    if err != nil {
        return err
    }
    err = binary.Write(writer, binary.BigEndian, uint32(len(data)))
    if err != nil {
        return err
    }
    _, err = writer.Write(data)
    return err
}
func readRequest(reader io.Reader) (*Request, error) {
    var length uint32
    err := binary.Read(reader, binary.BigEndian, &length)
    if err != nil {
        return nil, err
    }
    data := make([]byte, length)
    _, err = io.ReadFull(reader, data)
    if err != nil {
        return nil, err
    }
    var req Request
    err = json.Unmarshal(data, &req)
    return &req, err
}
func writeResponse(writer io.Writer, resp *Response) error {
    data, err := json.Marshal(resp)
    if err != nil {
        return err
    }
    err = binary.Write(writer, binary.BigEndian, uint32(len(data)))
    if err != nil {
        return err
    }
    _, err = writer.Write(data)
    return err
}
func readResponse(reader io.Reader) (*Response, error) {
    var length uint32
    err := binary.Read(reader, binary.BigEndian, &length)
    if err != nil {
        return nil, err
    }
    data := make([]byte, length)
    _, err = io.ReadFull(reader, data)
    if err != nil {
        return nil, err
    }
    var resp Response
    err = json.Unmarshal(data, &resp)
    return &resp, err
}