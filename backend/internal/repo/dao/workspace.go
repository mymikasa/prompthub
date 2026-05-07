package dao

import (
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

func (d *WorkspaceDAO) Create(ws *domain.Workspace) error {
	return d.db.Create(ws).Error
}

func (d *WorkspaceDAO) AddMember(member *domain.WorkspaceMember) error {
	return d.db.Create(member).Error
}

func (d *WorkspaceDAO) FindByUserID(userID int64) (*domain.Workspace, error) {
	var member domain.WorkspaceMember
	if err := d.db.Where("user_id = ? AND role = ?", userID, "owner").First(&member).Error; err != nil {
		return nil, err
	}
	var ws domain.Workspace
	if err := d.db.First(&ws, member.WorkspaceID).Error; err != nil {
		return nil, err
	}
	return &ws, nil
}
