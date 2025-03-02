package db

import (
	"time"

	"gorm.io/gorm"
)

type URL struct {
	ID        uint      `gorm:"primaryKey"`
	ShortCode string    `gorm:"uniqueIndex;not null"`
	OriginalURL string  `gorm:"not null"`
	CreatedAt time.Time
	ExpiresAt *time.Time
	Clicks    uint      `gorm:"default:0"`
	UserID    string    // Optional: if you want to track URL ownership
}

var DB *gorm.DB

func InitDB(dbPath string) error {
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return err
	}
	
	// Auto migrate the schema
	err = DB.AutoMigrate(&URL{})
	return err
}
