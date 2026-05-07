package domain

import "time"

type ProviderConfig struct {
	ID              int64     `json:"id"`
	WorkspaceID     int64     `json:"workspaceId"`
	ProviderType    string    `json:"providerType"`
	BaseURL         string    `json:"baseUrl"`
	EncryptedAPIKey string    `json:"-"`
	DefaultModel    string    `json:"defaultModel"`
	HasAPIKey       bool      `json:"hasApiKey"`
	CreatedBy       int64     `json:"createdBy"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}
