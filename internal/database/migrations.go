package database

import (
	"fmt"
	"github.com/souvik150/Fampay-Backend-Assignment/internal/models"
	"gorm.io/gorm"
	"log"
)

func RunMigrations(db *gorm.DB) {
	log.Println("Running Migrations")

	err := db.AutoMigrate(&models.User{}, &models.RefreshToken{}, &models.Video{})
	if err != nil {
		fmt.Println("Migration error")
		return
	}

	log.Println("🚀 Migrations completed")
}
