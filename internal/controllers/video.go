package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/souvik150/Fampay-Backend-Assignment/internal/services"
	"github.com/souvik150/Fampay-Backend-Assignment/internal/worker"
	"net/http"
)

func FetchAndStoreVideos(c *fiber.Ctx) error {
	topic := c.Params("topic")
	worker.StartVideoWorker(topic)

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Fetching videos from backend. The process is running in the background."})
}

func GetVideos(c *fiber.Ctx) error {
	topic := c.Query("topic", "")
	page := c.Query("page", "1")
	limit := c.Query("limit", "10")

	videos, totalResults, err := services.GetSortedVideos(page, limit, topic)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "totalResults": totalResults, "results": len(videos), "videos": videos})
}
