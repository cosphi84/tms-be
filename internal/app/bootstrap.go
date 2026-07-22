package app

import (
	"github.com/casbin/casbin/v3"
	"gorm.io/gorm"
)

type App struct {
	Enforcer *casbin.Enforcer
}

func TmsApp(db *gorm.DB) *App {
	return &App{}

	// casbin init
}
