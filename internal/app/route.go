package app

import (
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (app *Config) SetupRouter(
	router *gin.Engine,
	authMW gin.HandlerFunc,
	casbinMW gin.HandlerFunc,
) {
	corsOrigins := os.Getenv("APP_CORS_ORIGINS")
	if corsOrigins == "" {
		panic("APP_CORS_ORIGINS not setup yet. e.g: http://localhost:3000")
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Split(corsOrigins, ","),
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Request-ID"},
		AllowCredentials: true, // wajib true — tanpa ini cookie refresh_token gak akan pernah ke-set di browser FE
	}))

	api := router.Group("/")

	// ---- Public routes (gak butuh login) ----
	api.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"name":    app.AppName,
			"version": app.AppVersion,
		})
	})
	app.AuthHandler.RegisterPublicRoutes(api)
	app.OfficeHandler.RegisterPublicRoutes(api)

	// ---- Authenticated routes (butuh identity, TAPI TANPA RBAC) ----
	authenticated := api.Group("")
	authenticated.Use(authMW)
	{
		app.AuthHandler.RegisterProtectedRoutes(authenticated)
	}

	// ---- Fully protected routes (wajib login + lolos authorization) ----
	protected := api.Group("")
	protected.Use(authMW, casbinMW)
	{
		app.UploaderHandler.RegisterProtectedRoutes(protected)
		app.UsersHandler.RegisterRoutes(protected)
		app.OfficeHandler.RegisterProtectedRoutes(protected)
	}

}
