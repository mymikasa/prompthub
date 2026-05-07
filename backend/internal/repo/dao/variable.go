package dao

import (
	"context"

	"github.com/mymikasa/prompthub/internal/repo/dao/model"
	"gorm.io/gorm"
)

type VariableDAO interface {
	FindByPromptID(ctx context.Context, promptID int64) ([]*model.PromptVariable, error)
	FindByPromptIDAndName(ctx context.Context, promptID int64, name string) (*model.PromptVariable, error)
	FindByID(ctx context.Context, id int64) (*model.PromptVariable, error)
	BatchCreate(ctx context.Context, items []*model.PromptVariable) error
	Update(ctx context.Context, m *model.PromptVariable) error
	DeleteNotIn(ctx context.Context, promptID int64, names []string) error
	DeleteByPromptID(ctx context.Context, promptID int64) error
}

type variableDAO struct {
	db *gorm.DB
}

func NewVariableDAO(db *gorm.DB) VariableDAO {
	return &variableDAO{db: db}
}

func (d *variableDAO) FindByPromptID(ctx context.Context, promptID int64) ([]*model.PromptVariable, error) {
	var items []*model.PromptVariable
	if err := d.db.WithContext(ctx).Where("prompt_id = ?", promptID).Order("name ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (d *variableDAO) FindByPromptIDAndName(ctx context.Context, promptID int64, name string) (*model.PromptVariable, error) {
	var m model.PromptVariable
	if err := d.db.WithContext(ctx).Where("prompt_id = ? AND name = ?", promptID, name).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (d *variableDAO) FindByID(ctx context.Context, id int64) (*model.PromptVariable, error) {
	var m model.PromptVariable
	if err := d.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (d *variableDAO) BatchCreate(ctx context.Context, items []*model.PromptVariable) error {
	if len(items) == 0 {
		return nil
	}
	return d.db.WithContext(ctx).Create(&items).Error
}

func (d *variableDAO) Update(ctx context.Context, m *model.PromptVariable) error {
	return d.db.WithContext(ctx).Save(m).Error
}

func (d *variableDAO) DeleteNotIn(ctx context.Context, promptID int64, names []string) error {
	db := d.db.WithContext(ctx).Where("prompt_id = ?", promptID)
	if len(names) > 0 {
		db = db.Where("name NOT IN ?", names)
	}
	return db.Delete(&model.PromptVariable{}).Error
}

func (d *variableDAO) DeleteByPromptID(ctx context.Context, promptID int64) error {
	return d.db.WithContext(ctx).Where("prompt_id = ?", promptID).Delete(&model.PromptVariable{}).Error
}
