package main

import (
    "fmt"
    "github.com/valyala/fasthttp"
    "log"
    "os"
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
    // Security and environment variable compatibility
    publicIP, err := os.Hostname()
    if err != nil {
        log.Fatalln("Error getting public IP:", err)
    }

    fmt.Printf("[API Gateway] Public IP address: %s\n", publicIP)

    bindPort := os.Getenv("PORT")
    if bindPort == "" {
        bindPort = "8080"
    }

    f := &fasthttp.Server{
        Handler: routeHandler,
    }

    fmt.Printf("[API Gateway] Starting API gateway on port %s...\n", bindPort)
    if err := f.ListenAndServe(":" + bindPort); err != nil {
        log.Fatalln("[API Gateway] Error starting server:", err)
    }
}