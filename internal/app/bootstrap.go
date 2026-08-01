package app

import (
	"os"
	"tms-be/internal/auth"

	"tms-be/internal/casbin"
	"tms-be/internal/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Config struct {
	AppName    string
	AppVersion string
	DB         *gorm.DB

	// dependency yang dibutuhkan buat build middleware di main.go
	JWTSecret      string
	CasbinSvc      *casbin.Service
	AuthenticateMw gin.HandlerFunc
	AuthorizeMw    gin.HandlerFunc

	// per-module handler
	AuthHandler *auth.Handler
	//UsersHandler  *users.Handler
	//OfficeHandler *offices.Handler
}

func TmsAppBootstrap() *Config {
	app := &Config{
		AppName:    "tms",
		AppVersion: "0.1.0",
		JWTSecret:  os.Getenv("APP_JWT_SECRET"),
	}

	if app.JWTSecret == "" {
		panic("APP_JWT_SECRET not setup yet.")
	}

	// Init database
	db, err := database.Connect()
	if err != nil {
		panic(err)
	}
	app.DB = db

	// Init casbin enforcer (dipakai authorization middleware di main.go)
	casbinModelPath := os.Getenv("APP_CASBIN_MODEL_PATH")
	if casbinModelPath == "" {
		casbinModelPath = "conf/casbin_model.conf"
	}
	enforcer, err := casbin.NewEnforcer(db, casbinModelPath)
	if err != nil {
		panic(err)
	}
	app.CasbinSvc = casbin.New(enforcer)

	// middleware
	app.AuthenticateMw = auth.Authenticate()
	app.AuthorizeMw = casbin.Authorize(app.CasbinSvc)

	// ---- Wiring tiap module ----
	app.AuthHandler = auth.NewHandler(db)
	//app.UsersHandler = users.NewHandler(db)
	//app.OfficeHandler = offices.NewHandler(db)

	return app
}
