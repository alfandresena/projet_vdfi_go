package models

import (
	"time"

	"gorm.io/gorm"
)

type Ticket struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"uniqueIndex:event_user" json:"user_id"` // Un utilisateur ne peut avoir qu'un seul ticket par événement
	EventID   uint           `gorm:"index:event_user" json:"event_id"`
	Event     Event          `gorm:"foreignKey:EventID" json:"event"`
	User      User           `gorm:"foreignKey:UserID" json:"user"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
