package app

import (
	"tms-be/internal/database"

	"gorm.io/gorm"
)

type AppConfig struct {
	AppName    string
	AppVersion string
	DB         *gorm.DB

	UploaderHandler *uploader.UploaderHandler
}

func NewTmsApp() *AppConfig {
	app := &AppConfig{
		AppName:    "Tools Management System",
		AppVersion: "1.0.0",
	}

	// Connect to the database
	db, err := database.Connect()
	if err != nil {
		panic(err)
	}
	// Assign the database connection to the AppConfig struct
	app.DB = db

	return app
}
