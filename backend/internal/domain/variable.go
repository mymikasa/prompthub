package domain

import "time"

type PromptVariable struct {
	ID           int64     `json:"id"`
	PromptID     int64     `json:"promptId"`
	Name         string    `json:"name"`
	Label        string    `json:"label"`
	Description  string    `json:"description"`
	Required     bool      `json:"required"`
	DefaultValue string    `json:"defaultValue"`
	ExampleValue string    `json:"exampleValue"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
