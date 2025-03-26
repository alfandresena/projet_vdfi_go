package initializers

import "projet_vdfi/models"

func SyncDatabase() {
	// Sync database;
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Event{})
	DB.AutoMigrate(&models.Lyric{})
}
