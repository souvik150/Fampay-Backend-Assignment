package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/souvik150/Fampay-Backend-Assignment/internal/controllers"
)

func VideoRoutes(group fiber.Router) {
	videoGroup := group.Group("/video")

	// These are authenticated routes
	//videoGroup.Get("/", middleware.TokenValidation, controllers.GetVideos)
	//videoGroup.Post("/:topic", middleware.TokenValidation, controllers.FetchAndStoreVideos)

	// For the purpose of this task user without auth is also allowed
	videoGroup.Get("/", controllers.GetVideos)
	videoGroup.Post("/:topic", controllers.FetchAndStoreVideos)
}
