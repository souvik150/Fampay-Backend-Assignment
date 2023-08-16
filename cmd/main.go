package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/souvik150/Fampay-Backend-Assignment/config"
	"github.com/souvik150/Fampay-Backend-Assignment/internal/database"
	"github.com/souvik150/Fampay-Backend-Assignment/internal/routes"
	"log"
)

func main() {
	app := fiber.New()

	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatalln("Failed to load environment variables! \n", err.Error())
	}
	database.ConnectDB(&config)
	database.RunMigrations(database.DB)

	app.Use(logger.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins:     config.ClientOrigin,
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowMethods:     "GET, POST, PATCH, DELETE",
		AllowCredentials: true,
	}))

	apiGroup := app.Group("/v1")

	routes.AuthRoutes(apiGroup)

	apiGroup.Get("/healthcheck", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "Welcome to YouTube Video Fetcher API",
		})
	})

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Route not found",
		})
	})

	log.Fatal(app.Listen(config.Port))
}
