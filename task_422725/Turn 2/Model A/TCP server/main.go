package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

type Request struct {
	Method string `json:"method"`
	Path   string `json:"path"`
	Body   map[string]interface{} `json:"body,omitempty"`
}

type Response struct {
	Status int    `json:"status"`
	Message string `json:"message"`
	Body    map[string]interface{} `json:"body,omitempty"`
}

const (
	maxRequestSize = 1024 * 1024 // 1MB limit
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
	conn.SetDeadline(time.Now().Add(5 * time.Second))

	reader := bufio.NewReader(conn)
	for {
		data, err := readRequest(reader)
		if err != nil {
			if err == io.EOF {
				log.Println("Connection closed")
				return
			}
			log.Printf("Error reading request: %v", err)
			sendResponse(conn, 500, "Internal Server Error")
			continue
		}

		req := &Request{}
		if err := json.Unmarshal(data, req); err != nil {
			log.Printf("Error parsing request: %v", err)
			sendResponse(conn, 400, "Bad Request")
			continue
		}

		log.Printf("Received request: %+v", req)
		resp := handleRequest(req)
		sendResponse(conn, resp.Status, resp.Message, resp.Body)
	}
}

func readRequest(reader *bufio.Reader) ([]byte, error) {
	length, err := reader.ReadBytes(':')
	if err != nil {
		return nil, err
	}

	size, err := fmt.Sscanf(string(length), "%d", &length)
	if size != 1 || length > maxRequestSize {
		return nil, fmt.Errorf("invalid request size: %d", length)
	}

	return reader.ReadBytes('\n')
}

//Request handeling.

func handleRequest(req *Request) *Response {
	switch req.Method {
	case "GET":
		return getHandler(req)
	case "POST":
		return postHandler(req)
	default:
		return &Response{
			Status: 405,
			Message: "Method Not Allowed",
		}
	}
}

func getHandler(req *Request) *Response {
	return &Response{
		Status: 200,
		Message: "OK",
		Body: map[string]interface{}{
			"path": req.Path,
		},
	}
}

func postHandler(req *Request) *Response {
	return &Response{
		Status: 201,
		Message: "Created",
		Body:    req.Body,
	}
}