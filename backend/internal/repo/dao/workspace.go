package dao

import (
	"context"

	"github.com/mymikasa/prompthub/internal/repo/dao/model"
	"gorm.io/gorm"
)

type WorkspaceDAO struct {
	db *gorm.DB
}

func NewWorkspaceDAO(db *gorm.DB) *WorkspaceDAO {
	return &WorkspaceDAO{db: db}
}

func (d *WorkspaceDAO) Create(ctx context.Context, m *model.Workspace) error {
	return d.db.WithContext(ctx).Create(m).Error
}

func (d *WorkspaceDAO) AddMember(ctx context.Context, m *model.WorkspaceMember) error {
	return d.db.WithContext(ctx).Create(m).Error
}

func (d *WorkspaceDAO) FindMemberByUserID(ctx context.Context, userID int64) (*model.WorkspaceMember, error) {
	var m model.WorkspaceMember
	if err := d.db.WithContext(ctx).Where("user_id = ? AND role = ?", userID, "owner").First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (d *WorkspaceDAO) FindByID(ctx context.Context, id int64) (*model.Workspace, error) {
	var m model.Workspace
	if err := d.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}
