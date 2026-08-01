package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"tms-be/internal/app"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Get Port Configuration
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on OS environment variables")
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		panic("APP_PORT not setup yet.")
	}

	// Get Mode Configuration
	mode := os.Getenv("APP_MODE")
	if mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Bootstrap the Application
	// (nama fungsi harus sama persis dengan yang didefinisikan di bootstrap.go)
	cfg := app.TmsAppBootstrap()

	// Init Router
	route := gin.Default()

	// Set Trusted Proxies
	trustedProxies := os.Getenv("APP_TRUSTED_PROXIES")
	var proxies []string
	if trustedProxies != "" {
		proxies = strings.Split(trustedProxies, ",")
	}
	if err := route.SetTrustedProxies(proxies); err != nil {
		log.Fatal(err)
	}

	// Setup Router and bind to the application
	cfg.SetupRouter(route, cfg.AuthenticateMw, cfg.AuthorizeMw)

	// Start the server
	if err := route.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatal(err)
	}

}
