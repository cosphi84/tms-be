package main

import (
	"log"
	"os"
	"tms-be/internal/casbin"
	"tms-be/internal/database"
	"tms-be/internal/offices"
	"tms-be/internal/users"

	gormAdapter "github.com/casbin/gorm-adapter/v3"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("no .env file found, relying on OS environment variables")
		return
	}
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	if err := db.AutoMigrate(&gormAdapter.CasbinRule{}); err != nil {
		log.Fatal(err)
	}

	log.Println("HasTable:", db.Migrator().HasTable(&gormAdapter.CasbinRule{}))

	casbinModelPath := os.Getenv("APP_CASBIN_MODEL_PATH")
	if casbinModelPath == "" {
		casbinModelPath = "conf/casbin_model.conf"
	}
	enforcer, err := casbin.NewEnforcer(db, casbinModelPath)
	if err != nil {
		panic(err)
	}
	casbinSvc := casbin.New(enforcer)

	// Perform seeding operations
	hq, err := offices.Seed(db)
	if err != nil {
		log.Fatal(err)
	}
	err = users.Seed(db, casbinSvc, hq)
	if err != nil {
		log.Fatal(err)
	}

}
