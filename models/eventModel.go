package models

import (
	"time"

	"gorm.io/gorm"
)

type Event struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Title          string         `gorm:"not null" json:"title"`
	Description    string         `json:"description"`
	Location       string         `json:"location"`
	StartTime      time.Time      `json:"start_time"`
	EndTime        time.Time      `json:"end_time"`
	LiveStreamLink string         `json:"live_stream_link"` // Lien Zoom ou autre plateforme
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}
