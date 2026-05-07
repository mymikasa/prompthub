package repo

import (
	"context"
	"encoding/json"

	"github.com/mymikasa/prompthub/internal/domain"
	"github.com/mymikasa/prompthub/internal/repo/dao"
	"github.com/mymikasa/prompthub/internal/repo/dao/model"
)

type TestCaseRepo interface {
	Create(ctx context.Context, tc *domain.TestCase) error
	FindByID(ctx context.Context, id int64) (*domain.TestCase, error)
	FindByPromptID(ctx context.Context, promptID int64) ([]*domain.TestCase, error)
	Update(ctx context.Context, tc *domain.TestCase) error
	Delete(ctx context.Context, id int64) error
}

type testCaseRepo struct {
	dao dao.TestCaseDAO
}

func NewTestCaseRepo(d dao.TestCaseDAO) TestCaseRepo {
	return &testCaseRepo{dao: d}
}

func (r *testCaseRepo) Create(ctx context.Context, tc *domain.TestCase) error {
	m, err := toModelTestCase(tc)
	if err != nil {
		return err
	}
	if err := r.dao.Create(ctx, m); err != nil {
		return err
	}
	tc.ID = m.ID
	tc.CreatedAt = m.CreatedAt
	tc.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *testCaseRepo) FindByID(ctx context.Context, id int64) (*domain.TestCase, error) {
	m, err := r.dao.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toDomainTestCase(m)
}

func (r *testCaseRepo) FindByPromptID(ctx context.Context, promptID int64) ([]*domain.TestCase, error) {
	items, err := r.dao.FindByPromptID(ctx, promptID)
	if err != nil {
		return nil, err
	}
	result := make([]*domain.TestCase, len(items))
	for i, m := range items {
		result[i], err = toDomainTestCase(m)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (r *testCaseRepo) Update(ctx context.Context, tc *domain.TestCase) error {
	m, err := toModelTestCase(tc)
	if err != nil {
		return err
	}
	return r.dao.Update(ctx, m)
}

func (r *testCaseRepo) Delete(ctx context.Context, id int64) error {
	return r.dao.Delete(ctx, id)
}

func toModelTestCase(tc *domain.TestCase) (*model.PromptTestCase, error) {
	var valuesJSON string
	if tc.VariableValues != nil {
		b, err := json.Marshal(tc.VariableValues)
		if err != nil {
			return nil, err
		}
		valuesJSON = string(b)
	}
	return &model.PromptTestCase{
		ID:               tc.ID,
		PromptID:         tc.PromptID,
		Name:             tc.Name,
		VariableValues:   valuesJSON,
		ExpectedBehavior: tc.ExpectedBehavior,
		ExpectedOutput:   tc.ExpectedOutput,
		CreatedBy:        tc.CreatedBy,
		CreatedAt:        tc.CreatedAt,
		UpdatedAt:        tc.UpdatedAt,
	}, nil
}

func toDomainTestCase(m *model.PromptTestCase) (*domain.TestCase, error) {
	var values map[string]string
	if m.VariableValues != "" {
		if err := json.Unmarshal([]byte(m.VariableValues), &values); err != nil {
			return nil, err
		}
	}
	return &domain.TestCase{
		ID:               m.ID,
		PromptID:         m.PromptID,
		Name:             m.Name,
		VariableValues:   values,
		ExpectedBehavior: m.ExpectedBehavior,
		ExpectedOutput:   m.ExpectedOutput,
		CreatedBy:        m.CreatedBy,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
	}, nil
}
