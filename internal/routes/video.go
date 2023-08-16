package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/souvik150/Fampay-Backend-Assignment/internal/controllers"
)

func VideoRoutes(group fiber.Router) {
	videoGroup := group.Group("/video")
	videoGroup.Get("/", controllers.GetVideos)
	videoGroup.Post("/:topic", controllers.FetchAndStoreVideos)
}
