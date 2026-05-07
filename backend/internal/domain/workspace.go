package domain

import "time"

type Workspace struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	OwnerID   int64     `gorm:"index;not null" json:"owner_id"`
	CreatedAt time.Time `gorm:"index" json:"created_at"`
	UpdatedAt time.Time `gorm:"index" json:"updated_at"`
}

type WorkspaceMember struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	WorkspaceID int64     `gorm:"uniqueIndex:idx_ws_user;index;not null" json:"workspace_id"`
	UserID      int64     `gorm:"uniqueIndex:idx_ws_user;index;not null" json:"user_id"`
	Role        string    `gorm:"type:varchar(20);not null" json:"role"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
