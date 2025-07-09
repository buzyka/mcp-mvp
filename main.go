package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/buzyka/mcp-mvp/infrastructure/config"
	"github.com/buzyka/mcp-mvp/infrastructure/shopware"
	"github.com/gin-gonic/gin"
	"github.com/subosito/gotenv"
)

func main() {
    if fileExists(".env") {
		_ = gotenv.Load(".env")
	}

    cfg, err := config.NewFromEnv()
    if err != nil {
        panic(fmt.Sprintf("failed to load config: %v", err))
    }

    // configure these as needed
    targetDomain := cfg.TargetDomain

    authToken := "Bearer " + cfg.AccessToken
	clientID := cfg.ClientAccessKeyID
	clientSecret := cfg.ClientSecret

    fmt.Println(fmt.Sprintf("%#v", cfg))

    r := gin.Default()

    // catch any method on /mcp/:chat_id/*proxyPath
    r.Any("/mcp/:chat_id/*proxyPath", func(c *gin.Context) {
        chatID := c.Param("chat_id")
		fmt.Println("Chat ID:", chatID)
        proxyPath := c.Param("proxyPath") // includes leading slash, e.g. "/foo/bar"
        
        // Build target URL
        targetURL, err := url.Parse(targetDomain)
        if err != nil {
            log.Printf("invalid targetDomain: %v", err)
            c.String(http.StatusInternalServerError, "configuration error")
            return
        }
        targetURL.Path = proxyPath
        targetURL.RawQuery = c.Request.URL.RawQuery

        // Create new request with same method and body
        req, err := http.NewRequest(c.Request.Method, targetURL.String(), c.Request.Body)
        if err != nil {
            c.String(http.StatusInternalServerError, "failed to create request")
            return
        }
        defer c.Request.Body.Close()

        // Copy original headers (minus Host)
        for k, vv := range c.Request.Header {
            for _, v := range vv {
                req.Header.Add(k, v)
            }
        }

        // Inject auth header (or you could append as query param instead)
        req.Header.Set("Authorization", authToken)

        // Create Oauth2 client
        // Note: this client will use the token from the client credentials config
        client := shopware.NewSwClientFromIntegration(c, clientID, clientSecret, targetDomain)

        // Perform the request
        resp, err := client.HttpClient.Do(req)
        if err != nil {
            c.String(http.StatusBadGateway, "upstream request failed: %v", err)
            return
        }
        defer resp.Body.Close()

        // Copy upstream status
        c.Status(resp.StatusCode)
        // Copy upstream headers
        for k, vv := range resp.Header {
            for _, v := range vv {
                c.Writer.Header().Add(k, v)
            }
        }

        // Stream the body
        if _, err := io.Copy(c.Writer, resp.Body); err != nil {
            log.Printf("error streaming response: %v", err)
        }
    })

    // start server
    if err := r.Run(":8080"); err != nil {
        log.Fatalf("server failed: %v", err)
    }
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}
