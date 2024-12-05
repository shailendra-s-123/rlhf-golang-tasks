package main

import (
    "github.com/valyala/fasthttp"
)

func routeHandler(ctx *fasthttp.RequestCtx) {
    path := string(ctx.Path())
    switch path {
    case "/service1":
        // Route to service 1
        ctx.Reply([]byte("Response from Service 1"))
    case "/service2":
        // Route to service 2
        ctx.Reply([]byte("Response from Service 2"))
    default:
        ctx.ReplyStatusCode(fasthttp.StatusNotFound)
    }
}

func main() {
    f := &fasthttp.Server{
        Handler: routeHandler,
    }
    f.ListenAndServe(":8080")
}