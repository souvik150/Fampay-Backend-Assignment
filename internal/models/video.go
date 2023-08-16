package models

import (
	"time"
)

type Video struct {
	ID           uint64    `gorm:"primary_key" json:"id,omitempty"`
	Title        string    `gorm:"not null" json:"title,omitempty"`
	Description  string    `gorm:"not null" json:"description,omitempty"`
	PublishTime  time.Time `gorm:"not null" json:"publishTime,omitempty"`
	ThumbnailURL string    `gorm:"not null" json:"thumbnailURL,omitempty"`
}
