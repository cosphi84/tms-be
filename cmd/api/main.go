package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"tms-be/cmd/app"

	"github.com/gin-gonic/gin"
)

func main() {
	// Get Port Configuration
	port := os.Getenv("APP_PORT")
	if port == "" {
		panic("APP_PORT not setup yet.")
	}

	// Get Mode Configuration
	mode := os.Getenv("APP_MODE")
	if mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	//  Bootstrap the Application
	app := app.NewTmsApp()

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
	app.SetupRouter(route)

	// Start the server
	if err := route.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatal(err)
	}
}
