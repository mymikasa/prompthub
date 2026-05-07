package model

import "time"

type ProviderConfig struct {
	ID              int64     `gorm:"primaryKey;autoIncrement"`
	WorkspaceID     int64     `gorm:"index;not null"`
	ProviderType    string    `gorm:"type:varchar(50);not null"`
	BaseURL         string    `gorm:"type:varchar(500)"`
	EncryptedAPIKey string    `gorm:"type:text"`
	DefaultModel    string    `gorm:"type:varchar(100)"`
	CreatedBy       int64     `gorm:"index;not null"`
	CreatedAt       time.Time `gorm:"index"`
	UpdatedAt       time.Time `gorm:"index"`
}
