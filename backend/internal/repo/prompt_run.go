package repo

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/mymikasa/prompthub/internal/domain"
	"github.com/mymikasa/prompthub/internal/repo/dao"
	"github.com/mymikasa/prompthub/internal/repo/dao/model"
)

type RunFilter = dao.RunFilter

type PromptRunRepo interface {
	Create(ctx context.Context, run *domain.PromptRun) error
	FindByPromptID(ctx context.Context, promptID int64, page, pageSize int, filter RunFilter) ([]*domain.PromptRun, int64, error)
	FindByWorkspace(ctx context.Context, workspaceID int64, page, pageSize int, filter RunFilter) ([]*domain.PromptRun, int64, error)
}

type promptRunRepo struct {
	dao dao.PromptRunDAO
}

func NewPromptRunRepo(d dao.PromptRunDAO) PromptRunRepo {
	return &promptRunRepo{dao: d}
}

func (r *promptRunRepo) Create(ctx context.Context, run *domain.PromptRun) error {
	m, err := toModelPromptRun(run)
	if err != nil {
		return err
	}
	if err := r.dao.Create(ctx, m); err != nil {
		return err
	}
	run.ID = m.ID
	run.CreatedAt = m.CreatedAt
	return nil
}

func (r *promptRunRepo) FindByPromptID(ctx context.Context, promptID int64, page, pageSize int, filter RunFilter) ([]*domain.PromptRun, int64, error) {
	items, total, err := r.dao.FindByPromptID(ctx, promptID, page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}
	return toDomainRuns(items, total)
}

func (r *promptRunRepo) FindByWorkspace(ctx context.Context, workspaceID int64, page, pageSize int, filter RunFilter) ([]*domain.PromptRun, int64, error) {
	items, total, err := r.dao.FindByWorkspace(ctx, workspaceID, page, pageSize, filter)
	if err != nil {
		return nil, 0, err
	}
	return toDomainRuns(items, total)
}

func toDomainRuns(items []*model.PromptRun, total int64) ([]*domain.PromptRun, int64, error) {
	result := make([]*domain.PromptRun, len(items))
	for i, m := range items {
		r, err := toDomainPromptRun(m)
		if err != nil {
			return nil, 0, err
		}
		result[i] = r
	}
	return result, total, nil
}

func toModelPromptRun(run *domain.PromptRun) (*model.PromptRun, error) {
	varsJSON, err := json.Marshal(run.InputVariables)
	if err != nil {
		return nil, err
	}
	snapshotJSON, err := json.Marshal(run.RenderedPromptSnapshot)
	if err != nil {
		return nil, err
	}
	var usageJSON string
	if run.TokenUsage != nil {
		b, err := json.Marshal(run.TokenUsage)
		if err != nil {
			return nil, err
		}
		usageJSON = string(b)
	}
	var testCaseID sql.NullInt64
	if run.TestCaseID != nil {
		testCaseID = sql.NullInt64{Int64: *run.TestCaseID, Valid: true}
	}
	return &model.PromptRun{
		ID:                     run.ID,
		PromptID:               run.PromptID,
		PromptVersionID:        run.PromptVersionID,
		TestCaseID:             testCaseID,
		Provider:               run.Provider,
		Model:                  run.Model,
		InputVariables:         string(varsJSON),
		RenderedPromptSnapshot: string(snapshotJSON),
		OutputText:             run.OutputText,
		Latency:                run.Latency,
		TokenUsage:             usageJSON,
		ErrorMessage:           run.ErrorMessage,
		CreatedBy:              run.CreatedBy,
		CreatedAt:              run.CreatedAt,
	}, nil
}

func toDomainPromptRun(m *model.PromptRun) (*domain.PromptRun, error) {
	var inputVars map[string]string
	if m.InputVariables != "" {
		if err := json.Unmarshal([]byte(m.InputVariables), &inputVars); err != nil {
			return nil, err
		}
	}
	var snapshot any
	if m.RenderedPromptSnapshot != "" {
		if err := json.Unmarshal([]byte(m.RenderedPromptSnapshot), &snapshot); err != nil {
			return nil, err
		}
	}
	var usage *domain.TokenUsage
	if m.TokenUsage != "" {
		usage = &domain.TokenUsage{}
		if err := json.Unmarshal([]byte(m.TokenUsage), usage); err != nil {
			return nil, err
		}
	}
	var testCaseID *int64
	if m.TestCaseID.Valid {
		testCaseID = &m.TestCaseID.Int64
	}
	return &domain.PromptRun{
		ID:                     m.ID,
		PromptID:               m.PromptID,
		PromptVersionID:        m.PromptVersionID,
		TestCaseID:             testCaseID,
		Provider:               m.Provider,
		Model:                  m.Model,
		InputVariables:         inputVars,
		RenderedPromptSnapshot: snapshot,
		OutputText:             m.OutputText,
		Latency:                m.Latency,
		TokenUsage:             usage,
		ErrorMessage:           m.ErrorMessage,
		CreatedBy:              m.CreatedBy,
		CreatedAt:              m.CreatedAt,
	}, nil
}
