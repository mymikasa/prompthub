package repo

import (
	"context"

	"github.com/mymikasa/prompthub/internal/domain"
	"github.com/mymikasa/prompthub/internal/repo/dao"
	"github.com/mymikasa/prompthub/internal/repo/dao/model"
)

type PromptFilter = dao.PromptFilter

type PromptRepo interface {
	Create(ctx context.Context, p *domain.Prompt) error
	FindByID(ctx context.Context, id int64) (*domain.Prompt, error)
	FindByWorkspace(ctx context.Context, workspaceID int64, page, pageSize int, filter PromptFilter) ([]*domain.Prompt, int64, error)
	Update(ctx context.Context, p *domain.Prompt) error
}

type promptRepo struct {
	dao dao.PromptDAO
}

func NewPromptRepo(d dao.PromptDAO) PromptRepo {
	return &promptRepo{dao: d}
}

func (r *promptRepo) Create(ctx context.Context, p *domain.Prompt) error {
	m := toModelPrompt(p)
	if err := r.dao.Create(ctx, m); err != nil {
		return err
	}
	p.ID = m.ID
	p.CreatedAt = m.CreatedAt
	p.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *promptRepo) FindByID(ctx context.Context, id int64) (*domain.Prompt, error) {
	m, err := r.dao.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toDomainPrompt(m), nil
}

func (r *promptRepo) FindByWorkspace(ctx context.Context, workspaceID int64, page, pageSize int, filter PromptFilter) ([]*domain.Prompt, int64, error) {
	items, total, err := r.dao.FindByWorkspace(ctx, workspaceID, page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}
	result := make([]*domain.Prompt, len(items))
	for i, m := range items {
		result[i] = toDomainPrompt(m)
	}
	return result, total, nil
}

func (r *promptRepo) Update(ctx context.Context, p *domain.Prompt) error {
	return r.dao.Update(ctx, toModelPrompt(p))
}

func toModelPrompt(p *domain.Prompt) *model.Prompt {
	return &model.Prompt{
		ID:                 p.ID,
		WorkspaceID:        p.WorkspaceID,
		CreatedBy:          p.CreatedBy,
		Title:              p.Title,
		Slug:               p.Slug,
		Description:        p.Description,
		Body:               p.Body,
		MessageFormat:      p.MessageFormat,
		Visibility:         p.Visibility,
		Status:             p.Status,
		TargetProvider:     p.TargetProvider,
		TargetModel:        p.TargetModel,
		DefaultTemperature: p.DefaultTemperature,
		DefaultMaxTokens:   p.DefaultMaxTokens,
		UsageNotes:         p.UsageNotes,
		CreatedAt:          p.CreatedAt,
		UpdatedAt:          p.UpdatedAt,
	}
}

func toDomainPrompt(m *model.Prompt) *domain.Prompt {
	return &domain.Prompt{
		ID:                 m.ID,
		WorkspaceID:        m.WorkspaceID,
		CreatedBy:          m.CreatedBy,
		Title:              m.Title,
		Slug:               m.Slug,
		Description:        m.Description,
		Body:               m.Body,
		MessageFormat:      m.MessageFormat,
		Visibility:         m.Visibility,
		Status:             m.Status,
		TargetProvider:     m.TargetProvider,
		TargetModel:        m.TargetModel,
		DefaultTemperature: m.DefaultTemperature,
		DefaultMaxTokens:   m.DefaultMaxTokens,
		UsageNotes:         m.UsageNotes,
		CreatedAt:          m.CreatedAt,
		UpdatedAt:          m.UpdatedAt,
	}
}
