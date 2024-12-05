package main

import (
	"fmt"
	"net/http"
)

// Store valid API keys in a secure place in a real application
var validAPIKeys = map[string]bool{
	"your_valid_api_key": true,
}

func authenticate(r *http.Request) bool {
	// Extract the API key from the 'Authorization' header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return false // No API key provided
	}

	// Remove the "Bearer " prefix from the header value
	apiKey := authHeader[7:]

	// Check if the API key is valid
	return validAPIKeys[apiKey]
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
	// Validate the API key
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
