package dao

import (
	"context"

	"github.com/mymikasa/prompthub/internal/repo/dao/model"
	"gorm.io/gorm"
)

type PromptRunDAO interface {
	Create(ctx context.Context, m *model.PromptRun) error
	FindByPromptID(ctx context.Context, promptID int64, page, pageSize int) ([]*model.PromptRun, int64, error)
}

type promptRunDAO struct {
	db *gorm.DB
}

func NewPromptRunDAO(db *gorm.DB) PromptRunDAO {
	return &promptRunDAO{db: db}
}

func (d *promptRunDAO) Create(ctx context.Context, m *model.PromptRun) error {
	return d.db.WithContext(ctx).Create(m).Error
}

func (d *promptRunDAO) FindByPromptID(ctx context.Context, promptID int64, page, pageSize int) ([]*model.PromptRun, int64, error) {
	var total int64
	db := d.db.WithContext(ctx).Model(&model.PromptRun{}).Where("prompt_id = ?", promptID)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []*model.PromptRun
	offset := (page - 1) * pageSize
	if err := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}
