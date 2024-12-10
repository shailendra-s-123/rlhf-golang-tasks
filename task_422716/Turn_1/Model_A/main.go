package main

import (
        "fmt"
        "log"
        "net/http"
        _"os"
        "time"
)

//loggingMiddleware logs requests by recording the method, request path, and latency
func loggingMiddleware(handler http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                start := time.Now()
                defer func() {
                        latency := time.Since(start)
                        ms := latency.Milliseconds()
                        log.Printf("%s %s - %dms", r.Method, r.URL.Path, ms)
                }()

                handler.ServeHTTP(w, r)
        })
}

// errorMiddleware handlespanic recovery, logging and serving internal server errors (500)
func errorMiddleware(handler http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                defer func() {
                        if err := recover(); err != nil {
                                log.Printf("panic: %v", err)
                                http.Error(w, err.(error).Error(), http.StatusInternalServerError)
                        }
                }()
                handler.ServeHTTP(w, r)
        })
}

// indexHandler serves a simple "Hello, World!" response.
func indexHandler(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/" {
                http.NotFound(w, r)
                return
        }

        fmt.Fprintf(w, "Hello, World!")
}

func main() {
        fmt.Println("Server starting on :8080")
        //Create an empty router
        mux := http.NewServeMux()

        mux.HandleFunc("/", indexHandler)


        //custom middleware to handle errors and log requests
        handler := errorMiddleware(loggingMiddleware(mux))

        server := &http.Server{
                Addr:         ":8080",
                Handler:      handler,
                ReadTimeout:  10 * time.Second,
                WriteTimeout: 10 * time.Second,
                IdleTimeout:  120 * time.Second,
        }

        if err := server.ListenAndServe(); err != nil {
                log.Fatalf("Error starting server: %v", err)
        }
}