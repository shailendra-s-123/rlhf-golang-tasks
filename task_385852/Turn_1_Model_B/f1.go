package main

import (
	"fmt"
	"net/http"
)

var validAPIKeys = map[string]bool{
	"YOUR_VALID_API_KEY": true,
}

func authenticate(r *http.Request) bool {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return false // No API key provided
	}

	apiKey := authHeader
	return validAPIKeys[apiKey]
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
	if !authenticate(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	fmt.Fprintf(w, "Welcome to the protected resource!")
}

func main() {
	http.HandleFunc("/protected", protectedHandler)
	fmt.Println("Server listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
