package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/souvik150/Fampay-Backend-Assignment/internal/database"
	"github.com/souvik150/Fampay-Backend-Assignment/internal/models"
	"gorm.io/gorm"
	"net/http"
	"sort"
	"strconv"
	"time"
)

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

func GetSortedVideos(page, limit, topic string) ([]models.Video, int64, error) {
	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var videos []models.Video
	var totalResults int64

	dbQuery := database.DB
	if topic != "" {
		dbQuery = dbQuery.Where("topic = ?", topic)
	}

	results := dbQuery.Limit(intLimit).Offset(offset).Find(&videos)
	if results.Error != nil {
		return nil, 0, results.Error
	}

	sort.SliceStable(videos, func(i, j int) bool {
		return videos[i].PublishTime.After(videos[j].PublishTime)
	})

	countQuery := dbQuery.Model(&models.Video{})
	if topic != "" {
		countQuery = countQuery.Where("topic = ?", topic)
	}
	countResult := countQuery.Count(&totalResults)

	if countResult.Error != nil {
		return nil, 0, countResult.Error
	}

	return videos, totalResults, nil
}

func SaveVideo(video models.Video) error {
	// Check if a video with the same title already exists
	var existingVideo models.Video
	result := database.DB.Where("title = ?", video.Title).First(&existingVideo)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}

	// If a video with the same title already exists, return an error
	if result.RowsAffected > 0 {
		return errors.New("video with the same title already exists")
	}

	result = database.DB.Create(&video)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
