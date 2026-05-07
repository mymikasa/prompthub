package repo

import (
	"context"

	"github.com/mymikasa/prompthub/internal/domain"
	"github.com/mymikasa/prompthub/internal/repo/dao"
	"github.com/mymikasa/prompthub/internal/repo/dao/model"
)

type ProviderConfigRepo interface {
	ListByWorkspaceID(ctx context.Context, workspaceID int64) ([]domain.ProviderConfig, error)
	FindByID(ctx context.Context, id int64) (*domain.ProviderConfig, error)
	FindByWorkspaceAndType(ctx context.Context, workspaceID int64, providerType string) (*domain.ProviderConfig, error)
	Create(ctx context.Context, pc *domain.ProviderConfig) error
	Update(ctx context.Context, pc *domain.ProviderConfig) error
	Delete(ctx context.Context, id int64) error
}

type providerConfigRepo struct {
	dao dao.ProviderConfigDAO
}

func NewProviderConfigRepo(d dao.ProviderConfigDAO) ProviderConfigRepo {
	return &providerConfigRepo{dao: d}
}

func (r *providerConfigRepo) ListByWorkspaceID(ctx context.Context, workspaceID int64) ([]domain.ProviderConfig, error) {
	list, err := r.dao.ListByWorkspaceID(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	result := make([]domain.ProviderConfig, len(list))
	for i, m := range list {
		result[i] = *toDomainProviderConfig(&m)
	}
	return result, nil
}

func (r *providerConfigRepo) FindByID(ctx context.Context, id int64) (*domain.ProviderConfig, error) {
	m, err := r.dao.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toDomainProviderConfig(m), nil
}

func (r *providerConfigRepo) FindByWorkspaceAndType(ctx context.Context, workspaceID int64, providerType string) (*domain.ProviderConfig, error) {
	m, err := r.dao.FindByWorkspaceAndType(ctx, workspaceID, providerType)
	if err != nil {
		return nil, err
	}
	return toDomainProviderConfig(m), nil
}

func (r *providerConfigRepo) Create(ctx context.Context, pc *domain.ProviderConfig) error {
	m := toModelProviderConfig(pc)
	if err := r.dao.Create(ctx, m); err != nil {
		return err
	}
	pc.ID = m.ID
	pc.CreatedAt = m.CreatedAt
	pc.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *providerConfigRepo) Update(ctx context.Context, pc *domain.ProviderConfig) error {
	return r.dao.Update(ctx, toModelProviderConfig(pc))
}

func (r *providerConfigRepo) Delete(ctx context.Context, id int64) error {
	return r.dao.Delete(ctx, id)
}

func toModelProviderConfig(pc *domain.ProviderConfig) *model.ProviderConfig {
	return &model.ProviderConfig{
		ID:              pc.ID,
		WorkspaceID:     pc.WorkspaceID,
		ProviderType:    pc.ProviderType,
		Name:            pc.Name,
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
		Name:            m.Name,
		BaseURL:         m.BaseURL,
		EncryptedAPIKey: m.EncryptedAPIKey,
		HasAPIKey:       m.EncryptedAPIKey != "",
		DefaultModel:    m.DefaultModel,
		CreatedBy:       m.CreatedBy,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}
