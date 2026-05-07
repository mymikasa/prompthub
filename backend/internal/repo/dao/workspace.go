package dao

import (
	"context"

	"github.com/mymikasa/prompthub/internal/domain"
	"github.com/mymikasa/prompthub/internal/repo"
	"gorm.io/gorm"
)

type WorkspaceDAO struct {
	db *gorm.DB
}

func NewWorkspaceDAO(db *gorm.DB) repo.WorkspaceRepo {
	return &WorkspaceDAO{db: db}
}

func (d *WorkspaceDAO) Create(ctx context.Context, ws *domain.Workspace) error {
	return d.db.WithContext(ctx).Create(ws).Error
}

func (d *WorkspaceDAO) AddMember(ctx context.Context, member *domain.WorkspaceMember) error {
	return d.db.WithContext(ctx).Create(member).Error
}

func (d *WorkspaceDAO) FindByUserID(ctx context.Context, userID int64) (*domain.Workspace, error) {
	var member domain.WorkspaceMember
	if err := d.db.WithContext(ctx).Where("user_id = ? AND role = ?", userID, "owner").First(&member).Error; err != nil {
		return nil, err
	}
	var ws domain.Workspace
	if err := d.db.WithContext(ctx).First(&ws, member.WorkspaceID).Error; err != nil {
		return nil, err
	}
	return &ws, nil
}
