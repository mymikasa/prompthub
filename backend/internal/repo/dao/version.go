package dao

import (
	"context"

	"github.com/mymikasa/prompthub/internal/repo/dao/model"
	"gorm.io/gorm"
)

type VersionDAO interface {
	Create(ctx context.Context, m *model.PromptVersion) error
	FindByID(ctx context.Context, id int64) (*model.PromptVersion, error)
	FindByPromptID(ctx context.Context, promptID int64) ([]*model.PromptVersion, error)
	LatestVersion(ctx context.Context, promptID int64) (int, error)
}

type versionDAO struct {
	db *gorm.DB
}

func NewVersionDAO(db *gorm.DB) VersionDAO {
	return &versionDAO{db: db}
}

func (d *versionDAO) Create(ctx context.Context, m *model.PromptVersion) error {
	return d.db.WithContext(ctx).Create(m).Error
}

func (d *versionDAO) FindByID(ctx context.Context, id int64) (*model.PromptVersion, error) {
	var m model.PromptVersion
	if err := d.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (d *versionDAO) FindByPromptID(ctx context.Context, promptID int64) ([]*model.PromptVersion, error) {
	var items []*model.PromptVersion
	if err := d.db.WithContext(ctx).Where("prompt_id = ?", promptID).Order("version DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (d *versionDAO) LatestVersion(ctx context.Context, promptID int64) (int, error) {
	var m model.PromptVersion
	err := d.db.WithContext(ctx).Where("prompt_id = ?", promptID).Order("version DESC").First(&m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, nil
		}
		return 0, err
	}
	return m.Version, nil
}
