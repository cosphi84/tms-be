package main

import (
	"os"
	"tms-be/internal/database"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("APP_PORT")
	if port == "" {
		panic("APP_PORT not setup yet.")
	}

	mode := os.Getenv("APP_MODE")
	if mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	db, err := database.Connect()
	if err != nil {
		panic(err)
	}

}
