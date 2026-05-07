package dao

import (
	"context"

	"github.com/mymikasa/prompthub/internal/repo/dao/model"
	"gorm.io/gorm"
)

type PromptDAO interface {
	Create(ctx context.Context, m *model.Prompt) error
	FindByID(ctx context.Context, id int64) (*model.Prompt, error)
	FindByWorkspace(ctx context.Context, workspaceID int64, page, pageSize int) ([]*model.Prompt, int64, error)
	Update(ctx context.Context, m *model.Prompt) error
}

type promptDAO struct {
	db *gorm.DB
}

func NewPromptDAO(db *gorm.DB) PromptDAO {
	return &promptDAO{db: db}
}

func (d *promptDAO) Create(ctx context.Context, m *model.Prompt) error {
	return d.db.WithContext(ctx).Create(m).Error
}

func (d *promptDAO) FindByID(ctx context.Context, id int64) (*model.Prompt, error) {
	var m model.Prompt
	if err := d.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (d *promptDAO) FindByWorkspace(ctx context.Context, workspaceID int64, page, pageSize int) ([]*model.Prompt, int64, error) {
	var items []*model.Prompt
	var total int64

	db := d.db.WithContext(ctx).Model(&model.Prompt{}).Where("workspace_id = ?", workspaceID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Order("updated_at DESC").Offset(offset).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (d *promptDAO) Update(ctx context.Context, m *model.Prompt) error {
	return d.db.WithContext(ctx).Save(m).Error
}
