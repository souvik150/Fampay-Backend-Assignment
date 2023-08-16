package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/souvik150/Fampay-Backend-Assignment/internal/database"
	"github.com/souvik150/Fampay-Backend-Assignment/internal/models"
)

func TokenValidation(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "Authorization header missing"})
	}

	accessToken := authHeader[len("Bearer "):]
	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("your-access-token-secret-key"), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "Invalid or expired access token"})
	}

	userIDStr, ok := claims["userID"].(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "Invalid access token"})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "Invalid access token"})
	}

	// Check if the user exists in the database
	var user models.User
	result := database.DB.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "User not found"})
	}

	// Pass the authenticated user to the context
	c.Locals("userID", userID)

	return c.Next()
}
