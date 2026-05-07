package domain

import "time"

type ProviderConfig struct {
	ID              int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	WorkspaceID     int64     `gorm:"index;not null" json:"workspace_id"`
	ProviderType    string    `gorm:"type:varchar(50);not null" json:"provider_type"`
	BaseURL         string    `gorm:"type:varchar(500)" json:"base_url"`
	EncryptedAPIKey string    `gorm:"type:text" json:"-"`
	DefaultModel    string    `gorm:"type:varchar(100)" json:"default_model"`
	CreatedBy       int64     `gorm:"index;not null" json:"created_by"`
	CreatedAt       time.Time `gorm:"index" json:"created_at"`
	UpdatedAt       time.Time `gorm:"index" json:"updated_at"`
}
