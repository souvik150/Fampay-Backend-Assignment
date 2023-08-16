package database

import (
	"gorm.io/gorm"
	"log"
)

func RunMigrations(db *gorm.DB) {
	log.Println("Running Migrations")

	//db.AutoMigrate(&models.User{})
	//db.AutoMigrate(&models.Video{})

	log.Println("🚀 Migrations completed")
}
