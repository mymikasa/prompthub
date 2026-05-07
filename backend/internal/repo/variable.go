package repo

import (
	"context"

	"github.com/mymikasa/prompthub/internal/domain"
	"github.com/mymikasa/prompthub/internal/repo/dao"
	"github.com/mymikasa/prompthub/internal/repo/dao/model"
)

type VariableRepo interface {
	FindByPromptID(ctx context.Context, promptID int64) ([]*domain.PromptVariable, error)
	FindByID(ctx context.Context, id int64) (*domain.PromptVariable, error)
	FindByPromptIDAndName(ctx context.Context, promptID int64, name string) (*domain.PromptVariable, error)
	BatchCreate(ctx context.Context, items []*domain.PromptVariable) error
	Update(ctx context.Context, v *domain.PromptVariable) error
	DeleteNotIn(ctx context.Context, promptID int64, names []string) error
	DeleteByPromptID(ctx context.Context, promptID int64) error
}

type variableRepo struct {
	dao dao.VariableDAO
}

func NewVariableRepo(d dao.VariableDAO) VariableRepo {
	return &variableRepo{dao: d}
}

func (r *variableRepo) FindByPromptID(ctx context.Context, promptID int64) ([]*domain.PromptVariable, error) {
	items, err := r.dao.FindByPromptID(ctx, promptID)
	if err != nil {
		return nil, err
	}
	result := make([]*domain.PromptVariable, len(items))
	for i, m := range items {
		result[i] = toDomainVariable(m)
	}
	return result, nil
}

func (r *variableRepo) FindByID(ctx context.Context, id int64) (*domain.PromptVariable, error) {
	m, err := r.dao.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toDomainVariable(m), nil
}

func (r *variableRepo) FindByPromptIDAndName(ctx context.Context, promptID int64, name string) (*domain.PromptVariable, error) {
	m, err := r.dao.FindByPromptIDAndName(ctx, promptID, name)
	if err != nil {
		return nil, err
	}
	return toDomainVariable(m), nil
}

func (r *variableRepo) BatchCreate(ctx context.Context, items []*domain.PromptVariable) error {
	models := make([]*model.PromptVariable, len(items))
	for i, v := range items {
		models[i] = toModelVariable(v)
	}
	if err := r.dao.BatchCreate(ctx, models); err != nil {
		return err
	}
	for i, m := range models {
		items[i].ID = m.ID
		items[i].CreatedAt = m.CreatedAt
		items[i].UpdatedAt = m.UpdatedAt
	}
	return nil
}

func (r *variableRepo) Update(ctx context.Context, v *domain.PromptVariable) error {
	return r.dao.Update(ctx, toModelVariable(v))
}

func (r *variableRepo) DeleteNotIn(ctx context.Context, promptID int64, names []string) error {
	return r.dao.DeleteNotIn(ctx, promptID, names)
}

func (r *variableRepo) DeleteByPromptID(ctx context.Context, promptID int64) error {
	return r.dao.DeleteByPromptID(ctx, promptID)
}

func toModelVariable(v *domain.PromptVariable) *model.PromptVariable {
	return &model.PromptVariable{
		ID:           v.ID,
		PromptID:     v.PromptID,
		Name:         v.Name,
		Label:        v.Label,
		Description:  v.Description,
		Required:     v.Required,
		DefaultValue: v.DefaultValue,
		ExampleValue: v.ExampleValue,
		CreatedAt:    v.CreatedAt,
		UpdatedAt:    v.UpdatedAt,
	}
}

func toDomainVariable(m *model.PromptVariable) *domain.PromptVariable {
	return &domain.PromptVariable{
		ID:           m.ID,
		PromptID:     m.PromptID,
		Name:         m.Name,
		Label:        m.Label,
		Description:  m.Description,
		Required:     m.Required,
		DefaultValue: m.DefaultValue,
		ExampleValue: m.ExampleValue,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}
