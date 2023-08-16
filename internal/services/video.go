package services

import (
	"github.com/souvik150/Fampay-Backend-Assignment/internal/database"
	"github.com/souvik150/Fampay-Backend-Assignment/internal/models"
)

func SaveVideo(video models.Video) error {
	result := database.DB.Create(&video)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
