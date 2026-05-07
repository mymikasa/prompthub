package repo

import (
	"context"
	"encoding/json"

	"github.com/mymikasa/prompthub/internal/domain"
	"github.com/mymikasa/prompthub/internal/repo/dao"
	"github.com/mymikasa/prompthub/internal/repo/dao/model"
)

type VersionRepo interface {
	Create(ctx context.Context, v *domain.PromptVersion) error
	FindByID(ctx context.Context, id int64) (*domain.PromptVersion, error)
	FindByPromptID(ctx context.Context, promptID int64) ([]*domain.PromptVersion, error)
	LatestVersion(ctx context.Context, promptID int64) (int, error)
}

type versionRepo struct {
	dao dao.VersionDAO
}

func NewVersionRepo(d dao.VersionDAO) VersionRepo {
	return &versionRepo{dao: d}
}

func (r *versionRepo) Create(ctx context.Context, v *domain.PromptVersion) error {
	m, err := toModelVersion(v)
	if err != nil {
		return err
	}
	if err := r.dao.Create(ctx, m); err != nil {
		return err
	}
	v.ID = m.ID
	v.CreatedAt = m.CreatedAt
	return nil
}

func (r *versionRepo) FindByID(ctx context.Context, id int64) (*domain.PromptVersion, error) {
	m, err := r.dao.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toDomainVersion(m)
}

func (r *versionRepo) FindByPromptID(ctx context.Context, promptID int64) ([]*domain.PromptVersion, error) {
	items, err := r.dao.FindByPromptID(ctx, promptID)
	if err != nil {
		return nil, err
	}
	result := make([]*domain.PromptVersion, len(items))
	for i, m := range items {
		result[i], err = toDomainVersion(m)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (r *versionRepo) LatestVersion(ctx context.Context, promptID int64) (int, error) {
	return r.dao.LatestVersion(ctx, promptID)
}

func toModelVersion(v *domain.PromptVersion) (*model.PromptVersion, error) {
	snapshotBytes, err := json.Marshal(v.Snapshot)
	if err != nil {
		return nil, err
	}
	return &model.PromptVersion{
		ID:        v.ID,
		PromptID:  v.PromptID,
		Version:   v.VersionNumber,
		Snapshot:  string(snapshotBytes),
		CreatedBy: 0,
		CreatedAt: v.CreatedAt,
	}, nil
}

func toDomainVersion(m *model.PromptVersion) (*domain.PromptVersion, error) {
	var snapshot domain.VersionSnapshot
	if err := json.Unmarshal([]byte(m.Snapshot), &snapshot); err != nil {
		return nil, err
	}
	return &domain.PromptVersion{
		ID:            m.ID,
		PromptID:      m.PromptID,
		VersionNumber: m.Version,
		Snapshot:      snapshot,
		Author:        "",  // TODO: join user nickname
		CreatedAt:     m.CreatedAt,
	}, nil
}
