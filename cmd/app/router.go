package app

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (app *AppConfig) SetupRouter(router *gin.Engine) {
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization", "X-Request-ID"},
	}))

	api := router.Group("/")

	api.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"name":    app.AppName,
			"version": app.AppVersion,
		})
	})
}
