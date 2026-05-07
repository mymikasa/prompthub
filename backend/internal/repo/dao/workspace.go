package dao

import (
	"context"

	"github.com/mymikasa/prompthub/internal/repo/dao/model"
	"gorm.io/gorm"
)

type WorkspaceDAO interface {
	Create(ctx context.Context, m *model.Workspace) error
	AddMember(ctx context.Context, m *model.WorkspaceMember) error
	FindMemberByUserID(ctx context.Context, userID int64) (*model.WorkspaceMember, error)
	FindMemberByWorkspaceAndUser(ctx context.Context, workspaceID, userID int64) (*model.WorkspaceMember, error)
	FindByID(ctx context.Context, id int64) (*model.Workspace, error)
}

type workspaceDAO struct {
	db *gorm.DB
}

func NewWorkspaceDAO(db *gorm.DB) WorkspaceDAO {
	return &workspaceDAO{db: db}
}

func (d *workspaceDAO) Create(ctx context.Context, m *model.Workspace) error {
	return d.db.WithContext(ctx).Create(m).Error
}

func (d *workspaceDAO) AddMember(ctx context.Context, m *model.WorkspaceMember) error {
	return d.db.WithContext(ctx).Create(m).Error
}

func (d *workspaceDAO) FindMemberByUserID(ctx context.Context, userID int64) (*model.WorkspaceMember, error) {
	var m model.WorkspaceMember
	if err := d.db.WithContext(ctx).Where("user_id = ? AND role = ?", userID, "owner").First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (d *workspaceDAO) FindMemberByWorkspaceAndUser(ctx context.Context, workspaceID, userID int64) (*model.WorkspaceMember, error) {
	var m model.WorkspaceMember
	if err := d.db.WithContext(ctx).Where("workspace_id = ? AND user_id = ?", workspaceID, userID).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (d *workspaceDAO) FindByID(ctx context.Context, id int64) (*model.Workspace, error) {
	var m model.Workspace
	if err := d.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}
