package main

import (
    "fmt"
    "github.com/valyala/fasthttp"
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
    bindPort := os.Getenv("PORT")
    if bindPort == "" {
        bindPort = "8080"
        fmt.Println("Port not specified using environment variable, defaulting to 8080")
    }

    f := &fasthttp.Server{
        Handler: routeHandler,
    }

    fmt.Printf("Starting API gateway on port %s...\n", bindPort)
    f.ListenAndServe(":" + bindPort)
}