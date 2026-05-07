package domain

import "time"

type ProviderConfig struct {
	ID              int64     `json:"id"`
	WorkspaceID     int64     `json:"workspace_id"`
	ProviderType    string    `json:"provider_type"`
	BaseURL         string    `json:"base_url"`
	EncryptedAPIKey string    `json:"-"`
	DefaultModel    string    `json:"default_model"`
	CreatedBy       int64     `json:"created_by"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
