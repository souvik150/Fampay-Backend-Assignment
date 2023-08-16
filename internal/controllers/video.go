package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	initializers "github.com/souvik150/Fampay-Backend-Assignment/config"
	"github.com/souvik150/Fampay-Backend-Assignment/internal/database"
	"github.com/souvik150/Fampay-Backend-Assignment/internal/models"
	"github.com/souvik150/Fampay-Backend-Assignment/internal/services"
	"net/http"
	"strconv"
	"sync"
	"time"
)

func FetchAndStoreVideos(c *fiber.Ctx) error {
	config, _ := initializers.LoadConfig(".")
	apiKeys := config.APIKeysYT
	topic := c.Params("topic")
	totalDuration := 1 * time.Minute
	fetchInterval := 30 * time.Second

	stopFetching := time.After(totalDuration)
	videosChannel := make(chan []models.Video)
	var wg sync.WaitGroup

	for _, apiKey := range apiKeys {
		wg.Add(1)
		go func(apiKey string) {
			defer wg.Done()
			for {
				select {
				case <-stopFetching:
					return
				default:
					apiResponse, err := CallYouTubeAPI(apiKey, topic)
					if err != nil {
						fmt.Println("Error calling YouTube API:", err)
						continue
					}

					videosChannel <- apiResponse

					time.Sleep(fetchInterval)
				}
			}
		}(apiKey)
	}

	// Background goroutine to save videos
	go func() {
		for {
			select {
			case <-stopFetching:
				close(videosChannel) // Close the channel when fetching is done
				return
			case videos := <-videosChannel:
				for i := range videos {
					err := services.SaveVideo(videos[i])
					if err != nil {
						fmt.Println("Error storing videos:", err)
					}
				}
			}
		}
	}()

	wg.Wait()
	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Started fetching video from backend for the next 2 mins."})
}

func CallYouTubeAPI(apiKey string, searchQuery string) ([]models.Video, error) {
	currentTime := time.Now().UTC()
	publishedAfter := currentTime.Format(time.RFC3339)

	url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/search?key=%s&q=%s&type=video&order=date&part=snippet&publishedAfter=%s", apiKey, searchQuery, publishedAfter)

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("YouTube API request failed with status code: %d", response.StatusCode)
	}

	var apiResponse models.YouTubeAPIResponse
	err = json.NewDecoder(response.Body).Decode(&apiResponse)
	if err != nil {
		return nil, err
	}

	var videos []models.Video
	for _, item := range apiResponse.Items {
		publishTime, err := time.Parse(time.RFC3339, item.Snippet.PublishTime)
		if err != nil {
			return nil, err
		}

		video := models.Video{
			Topic:        searchQuery,
			Title:        item.Snippet.Title,
			Description:  item.Snippet.Description,
			PublishTime:  publishTime,
			ThumbnailURL: item.Snippet.Thumbnails.Default.URL,
		}
		videos = append(videos, video)
	}

	return videos, nil
}

func GetVideos(c *fiber.Ctx) error {
	var page = c.Query("page", "1")
	var limit = c.Query("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var videos []models.Video
	var totalResults int64

	results := database.DB.Limit(intLimit).Offset(offset).Find(&videos)
	if results.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": results.Error})
	}

	// Query for counting the total number of results
	countResult := database.DB.Model(&models.Video{}).Count(&totalResults)
	if countResult.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": countResult.Error.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "totalResults": totalResults, "results": len(videos), "videos": videos})
}
