package dao

import (
	"context"
	"time"

	"github.com/mymikasa/prompthub/internal/repo/dao/model"
	"gorm.io/gorm"
)

type RunFilter struct {
	Status    string // "success" or "failed"
	Model     string
	StartDate *time.Time
	EndDate   *time.Time
}

type PromptRunDAO interface {
	Create(ctx context.Context, m *model.PromptRun) error
	FindByPromptID(ctx context.Context, promptID int64, page, pageSize int, filter RunFilter) ([]*model.PromptRun, int64, error)
	FindByWorkspace(ctx context.Context, workspaceID int64, page, pageSize int, filter RunFilter) ([]*model.PromptRun, int64, error)
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

func (d *promptRunDAO) FindByPromptID(ctx context.Context, promptID int64, page, pageSize int, filter RunFilter) ([]*model.PromptRun, int64, error) {
	db := d.applyFilter(ctx, filter).Where("prompt_id = ?", promptID)
	return d.query(db, page, pageSize)
}

func (d *promptRunDAO) FindByWorkspace(ctx context.Context, workspaceID int64, page, pageSize int, filter RunFilter) ([]*model.PromptRun, int64, error) {
	db := d.applyFilter(ctx, filter).
		Joins("JOIN prompts ON prompts.id = prompt_runs.prompt_id").
		Where("prompts.workspace_id = ? AND prompts.deleted_at IS NULL", workspaceID)
	return d.query(db, page, pageSize)
}

func (d *promptRunDAO) applyFilter(ctx context.Context, filter RunFilter) *gorm.DB {
	db := d.db.WithContext(ctx).Model(&model.PromptRun{})
	if filter.Status == "success" {
		db = db.Where("error_message = ''")
	} else if filter.Status == "failed" {
		db = db.Where("error_message != ''")
	}
	if filter.Model != "" {
		db = db.Where("model = ?", filter.Model)
	}
	if filter.StartDate != nil {
		db = db.Where("created_at >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		db = db.Where("created_at <= ?", *filter.EndDate)
	}
	return db
}

func (d *promptRunDAO) query(db *gorm.DB, page, pageSize int) ([]*model.PromptRun, int64, error) {
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []*model.PromptRun
	offset := (page - 1) * pageSize
	if err := db.Order("prompt_runs.created_at DESC").Offset(offset).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}
