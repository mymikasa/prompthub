package service

import (
	"context"
	"errors"

	"github.com/mymikasa/prompthub/internal/domain"
	"github.com/mymikasa/prompthub/internal/repo"
	"gorm.io/gorm"
)

var ErrVariableNotFound = errors.New("variable not found")

type VariableService struct {
	variableRepo repo.VariableRepo
}

func NewVariableService(variableRepo repo.VariableRepo) *VariableService {
	return &VariableService{variableRepo: variableRepo}
}

type UpdateVariableReq struct {
	Label        *string `json:"label"`
	Description  *string `json:"description"`
	Required     *bool   `json:"required"`
	DefaultValue *string `json:"defaultValue"`
	ExampleValue *string `json:"exampleValue"`
}

func (s *VariableService) ListVariables(ctx context.Context, promptID int64) ([]*domain.PromptVariable, error) {
	return s.variableRepo.FindByPromptID(ctx, promptID)
}

func (s *VariableService) UpdateVariable(ctx context.Context, promptID, variableID int64, req *UpdateVariableReq) (*domain.PromptVariable, error) {
	v, err := s.variableRepo.FindByID(ctx, variableID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrVariableNotFound
		}
		return nil, err
	}
	if v.PromptID != promptID {
		return nil, ErrVariableNotFound
	}

	if req.Label != nil {
		v.Label = *req.Label
	}
	if req.Description != nil {
		v.Description = *req.Description
	}
	if req.Required != nil {
		v.Required = *req.Required
	}
	if req.DefaultValue != nil {
		v.DefaultValue = *req.DefaultValue
	}
	if req.ExampleValue != nil {
		v.ExampleValue = *req.ExampleValue
	}

	if err := s.variableRepo.Update(ctx, v); err != nil {
		return nil, err
	}
	return v, nil
}
