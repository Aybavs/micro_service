package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

// reverseProxy creates a reverse proxy handler that forwards incoming requests to the specified target service.
func reverseProxy(target string) gin.HandlerFunc {
    targetURL, err := url.Parse(target)
    if err != nil {
        log.Fatalf("Failed to parse target URL: %v", err)
    }
    proxy := httputil.NewSingleHostReverseProxy(targetURL)
    originalDirector := proxy.Director
    proxy.Director = func(req *http.Request) {
        originalDirector(req)
        req.Host = targetURL.Host
    }
    return func(c *gin.Context) {
        proxy.ServeHTTP(c.Writer, c.Request)
    }
}

func main() {
    r := gin.Default()

    r.Any("/users/*any", reverseProxy("http://user-service:8001"))
    r.Any("/products/*any", reverseProxy("http://product-service:8002"))

    r.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": "gateway up"})
    })

    r.Run(":8000")
}