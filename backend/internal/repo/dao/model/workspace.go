package model

import "time"

type Workspace struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"type:varchar(100);not null"`
	OwnerID   int64     `gorm:"index;not null"`
	CreatedAt time.Time `gorm:"index"`
	UpdatedAt time.Time `gorm:"index"`
}

type WorkspaceMember struct {
	ID          int64     `gorm:"primaryKey;autoIncrement"`
	WorkspaceID int64     `gorm:"uniqueIndex:idx_ws_user;index;not null"`
	UserID      int64     `gorm:"uniqueIndex:idx_ws_user;index;not null"`
	Role        string    `gorm:"type:varchar(20);not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
