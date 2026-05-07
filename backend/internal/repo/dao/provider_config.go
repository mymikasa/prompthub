package dao

import (
	"context"

	"github.com/mymikasa/prompthub/internal/repo/dao/model"
	"gorm.io/gorm"
)

type ProviderConfigDAO interface {
	ListByWorkspaceID(ctx context.Context, workspaceID int64) ([]model.ProviderConfig, error)
	FindByID(ctx context.Context, id int64) (*model.ProviderConfig, error)
	FindByWorkspaceAndType(ctx context.Context, workspaceID int64, providerType string) (*model.ProviderConfig, error)
	Create(ctx context.Context, m *model.ProviderConfig) error
	Update(ctx context.Context, m *model.ProviderConfig) error
	Delete(ctx context.Context, id int64) error
}

type providerConfigDAO struct {
	db *gorm.DB
}

func NewProviderConfigDAO(db *gorm.DB) ProviderConfigDAO {
	return &providerConfigDAO{db: db}
}

func (d *providerConfigDAO) ListByWorkspaceID(ctx context.Context, workspaceID int64) ([]model.ProviderConfig, error) {
	var list []model.ProviderConfig
	if err := d.db.WithContext(ctx).Where("workspace_id = ?", workspaceID).Order("id ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (d *providerConfigDAO) FindByID(ctx context.Context, id int64) (*model.ProviderConfig, error) {
	var m model.ProviderConfig
	if err := d.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (d *providerConfigDAO) FindByWorkspaceAndType(ctx context.Context, workspaceID int64, providerType string) (*model.ProviderConfig, error) {
	var m model.ProviderConfig
	if err := d.db.WithContext(ctx).Where("workspace_id = ? AND provider_type = ?", workspaceID, providerType).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (d *providerConfigDAO) Create(ctx context.Context, m *model.ProviderConfig) error {
	return d.db.WithContext(ctx).Create(m).Error
}

func (d *providerConfigDAO) Update(ctx context.Context, m *model.ProviderConfig) error {
	return d.db.WithContext(ctx).Save(m).Error
}

func (d *providerConfigDAO) Delete(ctx context.Context, id int64) error {
	return d.db.WithContext(ctx).Delete(&model.ProviderConfig{}, id).Error
}
