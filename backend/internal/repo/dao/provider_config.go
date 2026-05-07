package dao

import (
	"context"

	"github.com/mymikasa/prompthub/internal/repo/dao/model"
	"gorm.io/gorm"
)

type ProviderConfigDAO interface {
	FindByWorkspaceID(ctx context.Context, workspaceID int64) (*model.ProviderConfig, error)
	Upsert(ctx context.Context, m *model.ProviderConfig) error
}

type providerConfigDAO struct {
	db *gorm.DB
}

func NewProviderConfigDAO(db *gorm.DB) ProviderConfigDAO {
	return &providerConfigDAO{db: db}
}

func (d *providerConfigDAO) FindByWorkspaceID(ctx context.Context, workspaceID int64) (*model.ProviderConfig, error) {
	var m model.ProviderConfig
	if err := d.db.WithContext(ctx).Where("workspace_id = ?", workspaceID).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (d *providerConfigDAO) Upsert(ctx context.Context, m *model.ProviderConfig) error {
	var existing model.ProviderConfig
	err := d.db.WithContext(ctx).Where("workspace_id = ?", m.WorkspaceID).First(&existing).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return d.db.WithContext(ctx).Create(m).Error
		}
		return err
	}
	m.ID = existing.ID
	m.CreatedAt = existing.CreatedAt
	m.CreatedBy = existing.CreatedBy
	return d.db.WithContext(ctx).Save(m).Error
}
