package dao

import (
	"context"

	"github.com/mymikasa/prompthub/internal/repo/dao/model"
	"gorm.io/gorm"
)

type TestCaseDAO interface {
	Create(ctx context.Context, m *model.PromptTestCase) error
	FindByID(ctx context.Context, id int64) (*model.PromptTestCase, error)
	FindByPromptID(ctx context.Context, promptID int64) ([]*model.PromptTestCase, error)
	Update(ctx context.Context, m *model.PromptTestCase) error
	Delete(ctx context.Context, id int64) error
}

type testCaseDAO struct {
	db *gorm.DB
}

func NewTestCaseDAO(db *gorm.DB) TestCaseDAO {
	return &testCaseDAO{db: db}
}

func (d *testCaseDAO) Create(ctx context.Context, m *model.PromptTestCase) error {
	return d.db.WithContext(ctx).Create(m).Error
}

func (d *testCaseDAO) FindByID(ctx context.Context, id int64) (*model.PromptTestCase, error) {
	var m model.PromptTestCase
	if err := d.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (d *testCaseDAO) FindByPromptID(ctx context.Context, promptID int64) ([]*model.PromptTestCase, error) {
	var items []*model.PromptTestCase
	if err := d.db.WithContext(ctx).Where("prompt_id = ?", promptID).Order("created_at DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (d *testCaseDAO) Update(ctx context.Context, m *model.PromptTestCase) error {
	return d.db.WithContext(ctx).Save(m).Error
}

func (d *testCaseDAO) Delete(ctx context.Context, id int64) error {
	return d.db.WithContext(ctx).Delete(&model.PromptTestCase{}, id).Error
}
