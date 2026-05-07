package domain

import "time"

type Tag struct {
	ID          int64     `json:"id"`
	WorkspaceID int64     `json:"workspaceId"`
	Name        string    `json:"name"`
	CreatedAt   time.Time `json:"createdAt"`
}
