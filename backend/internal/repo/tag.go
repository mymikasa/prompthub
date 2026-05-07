package repo

import (
	"context"

	"github.com/mymikasa/prompthub/internal/domain"
	"github.com/mymikasa/prompthub/internal/repo/dao"
	"github.com/mymikasa/prompthub/internal/repo/dao/model"
)

type TagRepo interface {
	FindByWorkspace(ctx context.Context, workspaceID int64) ([]*domain.Tag, error)
	FindOrCreateByNames(ctx context.Context, workspaceID int64, names []string) ([]*domain.Tag, error)
	FindByPromptID(ctx context.Context, promptID int64) ([]*domain.Tag, error)
	ReplacePromptTags(ctx context.Context, promptID int64, tagIDs []int64) error
}

type tagRepo struct {
	dao dao.TagDAO
}

func NewTagRepo(d dao.TagDAO) TagRepo {
	return &tagRepo{dao: d}
}

func (r *tagRepo) FindByWorkspace(ctx context.Context, workspaceID int64) ([]*domain.Tag, error) {
	items, err := r.dao.FindByWorkspace(ctx, workspaceID)
	if err != nil {
		return nil, err
	}
	result := make([]*domain.Tag, len(items))
	for i, m := range items {
		result[i] = toDomainTag(m)
	}
	return result, nil
}

func (r *tagRepo) FindOrCreateByNames(ctx context.Context, workspaceID int64, names []string) ([]*domain.Tag, error) {
	items, err := r.dao.FindOrCreateByNames(ctx, workspaceID, names)
	if err != nil {
		return nil, err
	}
	result := make([]*domain.Tag, len(items))
	for i, m := range items {
		result[i] = toDomainTag(m)
	}
	return result, nil
}

func (r *tagRepo) FindByPromptID(ctx context.Context, promptID int64) ([]*domain.Tag, error) {
	items, err := r.dao.FindByPromptID(ctx, promptID)
	if err != nil {
		return nil, err
	}
	result := make([]*domain.Tag, len(items))
	for i, m := range items {
		result[i] = toDomainTag(m)
	}
	return result, nil
}

func (r *tagRepo) ReplacePromptTags(ctx context.Context, promptID int64, tagIDs []int64) error {
	return r.dao.ReplacePromptTags(ctx, promptID, tagIDs)
}

func toDomainTag(m *model.Tag) *domain.Tag {
	return &domain.Tag{
		ID:          m.ID,
		WorkspaceID: m.WorkspaceID,
		Name:        m.Name,
		CreatedAt:   m.CreatedAt,
	}
}
