package main

import (
	"fmt"
	"net/http"

	"github.com/Gkemhcs/k8s-admission-webhook/utils"
	"github.com/Gkemhcs/k8s-admission-webhook/webhook"
	"github.com/gin-gonic/gin"

	"github.com/spf13/pflag"
)

var (
	certFile string
	keyFile  string
	port     int
)

func init() {
	pflag.IntVar(&port, "port", 8080, "Webhook server port")
	pflag.StringVar(&certFile, "tls-cert-file", "/etc/webhook/certs/tls.crt", "Path to TLS certificate file")
	pflag.StringVar(&keyFile, "tls-key-file", "/etc/webhook/certs/tls.key", "Path to TLS key file")
	pflag.Parse()
}

func main() {
	logger := utils.NewLogger()

	// Create a new Gin router
	router := gin.Default()

	// Define a simple GET route
	router.GET("/", func(c *gin.Context) {
		logger.Info("Hello, TLS!")
		c.JSON(http.StatusOK, gin.H{"message": "Hello, TLS!"})
	})
	router.POST("/validate", webhook.ValidatePrivilegedContainer)
	
	if err := router.RunTLS(fmt.Sprintf("0.0.0.0:%d", port), certFile, keyFile); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
}
