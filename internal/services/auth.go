package services

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	initializers "github.com/souvik150/Fampay-Backend-Assignment/config"
	"github.com/souvik150/Fampay-Backend-Assignment/internal/database"
	"github.com/souvik150/Fampay-Backend-Assignment/internal/models"
	"github.com/souvik150/Fampay-Backend-Assignment/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func SignupUser(payload *models.RegisterUserSchema) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	otp, _ := utils.GenerateOTP(6)
	newUser := models.User{
		Username:     payload.Username,
		Email:        payload.Email,
		Password:     string(hashedPassword),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		ProfileImage: payload.ProfileImage,
		Verified:     false,
		Otp:          otp,
	}

	sub := "Account Verify"
	body := fmt.Sprintf("Please use your OTP: %s for account verification", otp)

	email, err := utils.SendEmail(payload.Email, sub, body)
	if err != nil {
		return err
	}
	fmt.Println(email)

	result := database.DB.Create(&newUser)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func LoginUser(payload *models.LoginUserSchema) (models.AuthResponse, error) {
	var user models.User
	result := database.DB.Where("username = ?", payload.Username).First(&user)
	if result.Error != nil {
		return models.AuthResponse{}, result.Error
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		return models.AuthResponse{}, err
	}

	if !user.Verified {
		return models.AuthResponse{}, nil
	}

	authResponse, err := GenerateAuthTokens(&user)
	if err != nil {
		return models.AuthResponse{}, err
	}

	return authResponse, nil
}

func VerifyOTP(userID uuid.UUID, otp string) error {
	var user models.User
	result := database.DB.Where("id = ? AND otp = ?", userID, otp).First(&user)
	if result.Error != nil {
		return result.Error
	}

	// Update user's verification status and clear OTP
	user.Verified = true
	user.Otp = ""
	result = database.DB.Save(&user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func GenerateAuthTokens(user *models.User) (models.AuthResponse, error) {
	accessToken, err := utils.GenerateAccessToken(user)
	if err != nil {
		return models.AuthResponse{}, err
	}

	refreshToken, err := utils.GenerateRefreshToken(user)
	if err != nil {
		return models.AuthResponse{}, err
	}

	// Store the refresh token in the database
	refreshTokenEntry := models.RefreshToken{
		UserID: user.ID,
		Token:  refreshToken,
	}
	result := database.DB.Create(&refreshTokenEntry)
	if result.Error != nil {
		return models.AuthResponse{}, result.Error
	}

	authResponse := models.AuthResponse{
		UserID:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Verified:     true,
	}

	return authResponse, nil
}

func ResendOTP(userID uuid.UUID) error {
	otp, _ := utils.GenerateOTP(6)

	result := database.DB.Model(&models.User{}).Where("id = ?", userID).Updates(models.User{
		Otp: otp,
	})
	if result.Error != nil {
		return result.Error
	}

	user, err := GetUserByID(userID)
	if err != nil {
		return err
	}

	if user.Verified == true {
		return nil
	}

	subject := "OTP Resent"
	body := fmt.Sprintf("Here is your OTP: %s", otp)

	email, err := utils.SendEmail(user.Email, subject, body)
	if err != nil {
		return err
	}
	fmt.Println(email)

	return nil
}

func RefreshAccessToken(refreshToken string) (models.AuthResponse, error) {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		return models.AuthResponse{}, err
	}

	// Check if the provided refresh token exists in the database
	var refreshTokenEntry models.RefreshToken
	result := database.DB.Where("token = ?", refreshToken).First(&refreshTokenEntry)
	if result.Error != nil {
		return models.AuthResponse{}, fiber.NewError(fiber.StatusUnauthorized, "Invalid refresh token")
	}

	// Parse and validate the access token
	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.RefreshTokenSecret), nil
	})
	if err != nil {
		return models.AuthResponse{}, fiber.NewError(fiber.StatusUnauthorized, "Invalid refresh token")
	}

	// Generate a new access token
	userIDStr, ok := claims["userID"].(string)
	if !ok {
		return models.AuthResponse{}, fiber.NewError(fiber.StatusUnauthorized, "Invalid refresh token")
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return models.AuthResponse{}, fiber.NewError(fiber.StatusUnauthorized, "Invalid refresh token")
	}

	var user models.User
	user.ID = userID
	accessToken, err := utils.GenerateAccessToken(&user)
	if err != nil {
		return models.AuthResponse{}, fiber.NewError(fiber.StatusInternalServerError, "Failed to generate access token")
	}

	authResponse := models.AuthResponse{
		UserID:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return authResponse, nil
}
