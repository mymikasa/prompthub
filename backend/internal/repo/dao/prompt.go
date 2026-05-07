package dao

import (
	"context"

	"github.com/mymikasa/prompthub/internal/repo/dao/model"
	"gorm.io/gorm"
)

type PromptFilter struct {
	Keyword   string
	Statuses  []string
	TagNames  []string
	Provider  string
	Model     string
}

type PromptDAO interface {
	Create(ctx context.Context, m *model.Prompt) error
	FindByID(ctx context.Context, id int64) (*model.Prompt, error)
	FindByWorkspace(ctx context.Context, workspaceID int64, page, pageSize int, filter PromptFilter) ([]*model.Prompt, int64, error)
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

func (d *promptDAO) FindByWorkspace(ctx context.Context, workspaceID int64, page, pageSize int, filter PromptFilter) ([]*model.Prompt, int64, error) {
	var items []*model.Prompt
	var total int64

	db := d.db.WithContext(ctx).Model(&model.Prompt{}).Where("workspace_id = ?", workspaceID)

	// default hide archived
	if len(filter.Statuses) == 0 {
		db = db.Where("status != ?", "archived")
	} else {
		db = db.Where("status IN ?", filter.Statuses)
	}

	if filter.Keyword != "" {
		like := "%" + filter.Keyword + "%"
		db = db.Where("title LIKE ? OR description LIKE ? OR body LIKE ?", like, like, like)
	}

	if filter.Provider != "" {
		db = db.Where("target_provider = ?", filter.Provider)
	}

	if filter.Model != "" {
		db = db.Where("target_model = ?", filter.Model)
	}

	if len(filter.TagNames) > 0 {
		db = db.Joins("JOIN prompt_tags ON prompt_tags.prompt_id = prompts.id").
			Joins("JOIN tags ON tags.id = prompt_tags.tag_id").
			Where("tags.name IN ?", filter.TagNames).
			Group("prompts.id")
	}

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
