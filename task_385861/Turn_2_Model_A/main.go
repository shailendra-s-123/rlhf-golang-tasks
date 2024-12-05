package main

import (
    "github.com/gin-gonic/gin"
    "log"
    "os"
)

func routeHandler(c *gin.Context) {
    path := c.Request.URL.Path
    switch path {
    case "/service1":
        c.String(200, "Response from Service 1")
    case "/service2":
        c.String(200, "Response from Service 2")
    default:
        c.String(404, "Not Found")
    }
}

func main() {
    env := os.Getenv("ENV")
    if env == "" {
        env = "development"
    }

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    router := gin.Default()
    router.GET("/service1", routeHandler)
    router.GET("/service2", routeHandler)

    log.Printf("[API Gateway] Starting in %s environment on port %s...\n", env, port)
    if err := router.Run(":" + port); err != nil {
        log.Fatalf("[API Gateway] Error starting server: %v\n", err)
    }
}