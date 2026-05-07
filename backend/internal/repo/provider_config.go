package repo

import (
	"context"

	"github.com/mymikasa/prompthub/internal/domain"
	"github.com/mymikasa/prompthub/internal/repo/dao"
	"github.com/mymikasa/prompthub/internal/repo/dao/model"
)

type ProviderConfigRepo interface {
	FindByWorkspaceID(ctx context.Context, workspaceID int64) (*domain.ProviderConfig, error)
	Save(ctx context.Context, pc *domain.ProviderConfig) error
}

type providerConfigRepo struct {
	dao dao.ProviderConfigDAO
}

func NewProviderConfigRepo(d dao.ProviderConfigDAO) ProviderConfigRepo {
	return &providerConfigRepo{dao: d}
}

func (r *providerConfigRepo) FindByWorkspaceID(ctx context.Context, workspaceID int64) (*domain.ProviderConfig, error) {
	m, err := r.dao.FindByWorkspaceID(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	return toDomainProviderConfig(m), nil
}

func (r *providerConfigRepo) Save(ctx context.Context, pc *domain.ProviderConfig) error {
	m := toModelProviderConfig(pc)
	if err := r.dao.Upsert(ctx, m); err != nil {
		return err
	}
	pc.ID = m.ID
	pc.CreatedAt = m.CreatedAt
	pc.UpdatedAt = m.UpdatedAt
	return nil
}

func toModelProviderConfig(pc *domain.ProviderConfig) *model.ProviderConfig {
	return &model.ProviderConfig{
		ID:              pc.ID,
		WorkspaceID:     pc.WorkspaceID,
		ProviderType:    pc.ProviderType,
		BaseURL:         pc.BaseURL,
		EncryptedAPIKey: pc.EncryptedAPIKey,
		DefaultModel:    pc.DefaultModel,
		CreatedBy:       pc.CreatedBy,
		CreatedAt:       pc.CreatedAt,
		UpdatedAt:       pc.UpdatedAt,
	}
}

func toDomainProviderConfig(m *model.ProviderConfig) *domain.ProviderConfig {
	return &domain.ProviderConfig{
		ID:              m.ID,
		WorkspaceID:     m.WorkspaceID,
		ProviderType:    m.ProviderType,
		BaseURL:         m.BaseURL,
		EncryptedAPIKey: m.EncryptedAPIKey,
		HasAPIKey:       m.EncryptedAPIKey != "",
		DefaultModel:    m.DefaultModel,
		CreatedBy:       m.CreatedBy,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}
