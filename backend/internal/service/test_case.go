package service

import (
	"context"
	"errors"

	"github.com/mymikasa/prompthub/internal/domain"
	"github.com/mymikasa/prompthub/internal/repo"
	"gorm.io/gorm"
)

var ErrTestCaseNotFound = errors.New("test case not found")

type TestCaseService struct {
	testCaseRepo repo.TestCaseRepo
	promptRepo   repo.PromptRepo
	variableRepo repo.VariableRepo
}

func NewTestCaseService(testCaseRepo repo.TestCaseRepo, promptRepo repo.PromptRepo, variableRepo repo.VariableRepo) *TestCaseService {
	return &TestCaseService{testCaseRepo: testCaseRepo, promptRepo: promptRepo, variableRepo: variableRepo}
}

type CreateTestCaseReq struct {
	Name             string            `json:"name" binding:"required,min=1,max=200"`
	VariableValues   map[string]string `json:"variableValues"`
	ExpectedBehavior string            `json:"expectedBehavior"`
	ExpectedOutput   string            `json:"expectedOutput"`
}

type UpdateTestCaseReq struct {
	Name             *string            `json:"name"`
	VariableValues   *map[string]string `json:"variableValues"`
	ExpectedBehavior *string            `json:"expectedBehavior"`
	ExpectedOutput   *string            `json:"expectedOutput"`
}

func (s *TestCaseService) List(ctx context.Context, actor *domain.Actor, promptID int64) ([]*domain.TestCase, error) {
	if _, err := s.getPrompt(ctx, actor, promptID); err != nil {
		return nil, err
	}
	return s.testCaseRepo.FindByPromptID(ctx, promptID)
}

func (s *TestCaseService) Create(ctx context.Context, actor *domain.Actor, promptID int64, req *CreateTestCaseReq) (*domain.TestCase, error) {
	p, err := s.getPrompt(ctx, actor, promptID)
	if err != nil {
		return nil, err
	}
	if !p.CanEdit(actor) {
		return nil, ErrNoPermission
	}

	if req.VariableValues != nil {
		if err := s.validateVariableValues(ctx, promptID, req.VariableValues); err != nil {
			return nil, err
		}
	}

	tc := &domain.TestCase{
		PromptID:         promptID,
		Name:             req.Name,
		VariableValues:   req.VariableValues,
		ExpectedBehavior: req.ExpectedBehavior,
		ExpectedOutput:   req.ExpectedOutput,
		CreatedBy:        actor.UserID,
	}
	if err := s.testCaseRepo.Create(ctx, tc); err != nil {
		return nil, err
	}
	return tc, nil
}

func (s *TestCaseService) Update(ctx context.Context, actor *domain.Actor, promptID, testCaseID int64, req *UpdateTestCaseReq) (*domain.TestCase, error) {
	p, err := s.getPrompt(ctx, actor, promptID)
	if err != nil {
		return nil, err
	}
	if !p.CanEdit(actor) {
		return nil, ErrNoPermission
	}

	tc, err := s.testCaseRepo.FindByID(ctx, testCaseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTestCaseNotFound
		}
		return nil, err
	}
	if tc.PromptID != promptID {
		return nil, ErrTestCaseNotFound
	}

	if req.Name != nil {
		tc.Name = *req.Name
	}
	if req.VariableValues != nil {
		if err := s.validateVariableValues(ctx, promptID, *req.VariableValues); err != nil {
			return nil, err
		}
		tc.VariableValues = *req.VariableValues
	}
	if req.ExpectedBehavior != nil {
		tc.ExpectedBehavior = *req.ExpectedBehavior
	}
	if req.ExpectedOutput != nil {
		tc.ExpectedOutput = *req.ExpectedOutput
	}

	if err := s.testCaseRepo.Update(ctx, tc); err != nil {
		return nil, err
	}
	return tc, nil
}

func (s *TestCaseService) Delete(ctx context.Context, actor *domain.Actor, promptID, testCaseID int64) error {
	p, err := s.getPrompt(ctx, actor, promptID)
	if err != nil {
		return err
	}
	if !p.CanEdit(actor) {
		return ErrNoPermission
	}

	tc, err := s.testCaseRepo.FindByID(ctx, testCaseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTestCaseNotFound
		}
		return err
	}
	if tc.PromptID != promptID {
		return ErrTestCaseNotFound
	}

	return s.testCaseRepo.Delete(ctx, testCaseID)
}

func (s *TestCaseService) getPrompt(ctx context.Context, actor *domain.Actor, promptID int64) (*domain.Prompt, error) {
	p, err := s.promptRepo.FindByID(ctx, promptID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPromptNotFound
		}
		return nil, err
	}
	if !p.CanView(actor.UserID) {
		return nil, ErrPromptNotFound
	}
	return p, nil
}

func (s *TestCaseService) validateVariableValues(ctx context.Context, promptID int64, values map[string]string) error {
	variables, err := s.variableRepo.FindByPromptID(ctx, promptID)
	if err != nil {
		return err
	}
	validNames := make(map[string]struct{}, len(variables))
	for _, v := range variables {
		validNames[v.Name] = struct{}{}
	}
	for k := range values {
		if _, ok := validNames[k]; !ok {
			return errors.New("unknown variable: " + k)
		}
	}
	return nil
}
