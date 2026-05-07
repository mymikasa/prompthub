package model

import "time"

type Tag struct {
	ID          int64     `gorm:"primaryKey;autoIncrement"`
	WorkspaceID int64     `gorm:"uniqueIndex:idx_ws_tag_name;index;not null"`
	Name        string    `gorm:"type:varchar(50);uniqueIndex:idx_ws_tag_name;not null"`
	CreatedAt   time.Time
}

type PromptTag struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	PromptID  int64     `gorm:"uniqueIndex:idx_prompt_tag;index;not null"`
	TagID     int64     `gorm:"uniqueIndex:idx_prompt_tag;index;not null"`
	CreatedAt time.Time
}
