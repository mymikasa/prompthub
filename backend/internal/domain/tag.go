package domain

import "time"

type Tag struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	WorkspaceID int64     `gorm:"uniqueIndex:idx_ws_tag_name;index;not null" json:"workspace_id"`
	Name        string    `gorm:"type:varchar(50);uniqueIndex:idx_ws_tag_name;not null" json:"name"`
	CreatedAt   time.Time `json:"created_at"`
}

type PromptTag struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	PromptID  int64     `gorm:"uniqueIndex:idx_prompt_tag;index;not null" json:"prompt_id"`
	TagID     int64     `gorm:"uniqueIndex:idx_prompt_tag;index;not null" json:"tag_id"`
	CreatedAt time.Time `json:"created_at"`
}
