package domain

import "time"

type Tag struct {
	ID          int64     `json:"id"`
	WorkspaceID int64     `json:"workspace_id"`
	Name        string    `json:"name"`
	CreatedAt   time.Time `json:"created_at"`
}

type PromptTag struct {
	ID        int64     `json:"id"`
	PromptID  int64     `json:"prompt_id"`
	TagID     int64     `json:"tag_id"`
	CreatedAt time.Time `json:"created_at"`
}
