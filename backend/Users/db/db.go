package db

import (
	"log"
	"os"

	"github.com/gajare/Fish-market/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set.")
	}
	gbd, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("gorm open error:", err)
	}
	//automatically migrate the User model
	if err := gbd.AutoMigrate(&models.User{}); err != nil {
		log.Fatal("gorm auto migrate error:", err)
	}

	DB = gbd
	log.Println("connected to postgres using gorm!")
}
