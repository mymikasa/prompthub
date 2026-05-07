package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           int64          `gorm:"primaryKey;autoIncrement"`
	Email        string         `gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash string         `gorm:"type:varchar(255);not null"`
	Nickname     string         `gorm:"type:varchar(100)"`
	CreatedAt    time.Time      `gorm:"index"`
	UpdatedAt    time.Time      `gorm:"index"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}
