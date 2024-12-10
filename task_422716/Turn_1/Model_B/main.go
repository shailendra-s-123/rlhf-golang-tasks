package main

import (
    "fmt"
    "log"
    "net/http"
)

// define a custom logger type for middleware
type logger struct {
    next http.Handler
}

// Wrap another handler and log each incoming request
func (l logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    log.Printf("Request method: %s, URL: %s\n", r.Method, r.URL)
    l.next.ServeHTTP(w, r)
}

// A generic error handler that logs errors
func errorHandler(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                log.Println("Error:", err)
                http.Error(w, err.Error(), http.StatusInternalServerError)
            }
        }()
        next.ServeHTTP(w, r)
    })
}

// Example handler function to simulate some functionality
func exampleHandler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path == "/error" {
        panic("Internal Server Error")
    }
    fmt.Fprintf(w, "Hello, %s!", r.URL.Path)
}

func main() {
    // Create the logger middleware
    logMiddleware := logger{next: http.HandlerFunc(exampleHandler)}

    // Wrap the logger with the error handler
    handler := errorHandler(logMiddleware)

    // Start the server
    log.Fatal(http.ListenAndServe(":8080", handler))
}