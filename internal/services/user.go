package services

import (
	"github.com/google/uuid"
	"github.com/souvik150/Fampay-Backend-Assignment/internal/database"
	"github.com/souvik150/Fampay-Backend-Assignment/internal/models"
)

func GetUserByID(userID uuid.UUID) (models.User, error) {
	var user models.User
	result := database.DB.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		return models.User{}, result.Error
	}
	return user, nil
}

func GetUserByEmail(email string) (models.User, error) {
	var user models.User
	result := database.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return models.User{}, result.Error
	}
	return user, nil
}

func GetUserByUsername(username string) (models.User, error) {
	var user models.User
	result := database.DB.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return models.User{}, result.Error
	}
	return user, nil
}
