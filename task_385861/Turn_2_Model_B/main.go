package main

import (
    "flag"
    "fmt"
    "github.com/valyala/fasthttp"
    "log"
    "os"
    "strconv"
)

var (
    port = flag.String("port", "8080", "Port to listen on")
)

func routeHandler(ctx *fasthttp.RequestCtx) {
    path := string(ctx.Path())
    switch path {
    case "/service1":
        ctx.Reply([]byte("Response from Service 1"))
    case "/service2":
        ctx.Reply([]byte("Response from Service 2"))
    default:
        ctx.ReplyStatusCode(fasthttp.StatusNotFound)
    }
}

func main() {
    flag.Parse()

    // Read port from environment variable if set
    if envPort := os.Getenv("PORT"); envPort != "" {
        portValue, err := strconv.Atoi(envPort)
        if err != nil {
            log.Fatalf("Invalid port environment variable: %v", err)
        }
        *port = fmt.Sprintf("%d", portValue)
    }

    f := &fasthttp.Server{
        Handler: routeHandler,
    }

    fmt.Printf("Starting API gateway on port %s...\n", *port)
    if err := f.ListenAndServe(":" + *port); err != nil {
        log.Fatalln("Error starting server:", err)
    }
}