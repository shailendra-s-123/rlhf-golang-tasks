package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/service1/", service1Handler)
	http.HandleFunc("/service2/", service2Handler)

	fmt.Println("API Gateway listening on port", port)
	http.ListenAndServe(":"+port, nil)
}

func service1Handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Routing to Service 1"))
}

func service2Handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Routing to Service 2"))
}