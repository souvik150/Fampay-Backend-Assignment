package models

import (
	"time"
)

type Video struct {
	ID           uint64    `gorm:"primary_key" json:"id,omitempty"`
	Topic        string    `gorm:"not null" json:"topic,omitempty"`
	Title        string    `gorm:"unique" json:"title,omitempty"`
	Description  string    `gorm:"not null" json:"description,omitempty"`
	PublishTime  time.Time `gorm:"not null" json:"publishTime,omitempty"`
	ThumbnailURL string    `gorm:"not null" json:"thumbnailURL,omitempty"`
}

type YouTubeAPIResponse struct {
	Items []struct {
		Snippet struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			PublishTime string `json:"publishedAt"`
			Thumbnails  struct {
				Default struct {
					URL string `json:"url"`
				} `json:"default"`
			} `json:"thumbnails"`
		} `json:"snippet"`
	} `json:"items"`
}
