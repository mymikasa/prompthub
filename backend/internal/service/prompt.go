package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/mymikasa/prompthub/internal/domain"
	"github.com/mymikasa/prompthub/internal/repo"
	"gorm.io/gorm"
)

var (
	ErrPromptNotFound  = errors.New("prompt not found")
	ErrNoPermission    = errors.New("no permission")
)

type PromptService struct {
	promptRepo repo.PromptRepo
}

func NewPromptService(promptRepo repo.PromptRepo) *PromptService {
	return &PromptService{promptRepo: promptRepo}
}

type CreatePromptReq struct {
	Title              string   `json:"title" binding:"required,min=1,max=200"`
	Slug               string   `json:"slug"`
	Description        string   `json:"description"`
	Body               string   `json:"body" binding:"required"`
	MessageFormat      string   `json:"messageFormat" binding:"required,oneof=single_text chat_messages"`
	Visibility         string   `json:"visibility" binding:"required,oneof=private workspace"`
	TargetProvider     string   `json:"targetProvider"`
	TargetModel        string   `json:"targetModel"`
	DefaultTemperature *float64 `json:"defaultTemperature"`
	DefaultMaxTokens   *int     `json:"defaultMaxTokens"`
	UsageNotes         string   `json:"usageNotes"`
}

type UpdatePromptReq struct {
	Title              *string  `json:"title"`
	Slug               *string  `json:"slug"`
	Description        *string  `json:"description"`
	Body               *string  `json:"body"`
	MessageFormat      *string  `json:"messageFormat"`
	Visibility         *string  `json:"visibility"`
	Status             *string  `json:"status"`
	TargetProvider     *string  `json:"targetProvider"`
	TargetModel        *string  `json:"targetModel"`
	DefaultTemperature *float64 `json:"defaultTemperature"`
	DefaultMaxTokens   *int     `json:"defaultMaxTokens"`
	UsageNotes         *string  `json:"usageNotes"`
}

func (s *PromptService) CreatePrompt(ctx context.Context, actor *domain.Actor, req *CreatePromptReq) (*domain.Prompt, error) {
	slug := req.Slug
	if slug == "" {
		slug = generateSlug(req.Title)
	}

	p := &domain.Prompt{
		WorkspaceID:        actor.WorkspaceID,
		CreatedBy:          actor.UserID,
		Title:              req.Title,
		Slug:               slug,
		Description:        req.Description,
		Body:               req.Body,
		MessageFormat:      req.MessageFormat,
		Visibility:         req.Visibility,
		Status:             "draft",
		TargetProvider:     req.TargetProvider,
		TargetModel:        req.TargetModel,
		DefaultTemperature: req.DefaultTemperature,
		DefaultMaxTokens:   req.DefaultMaxTokens,
		UsageNotes:         req.UsageNotes,
	}

	if err := s.promptRepo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *PromptService) GetPrompt(ctx context.Context, actor *domain.Actor, id int64) (*domain.Prompt, error) {
	p, err := s.promptRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPromptNotFound
		}
		return nil, err
	}

	if p.WorkspaceID != actor.WorkspaceID {
		return nil, ErrPromptNotFound
	}

	if !p.CanView(actor.UserID) {
		return nil, ErrPromptNotFound
	}

	return p, nil
}

func (s *PromptService) ListPrompts(ctx context.Context, actor *domain.Actor, page, pageSize int) (*domain.PromptList, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	items, total, err := s.promptRepo.FindByWorkspace(ctx, actor.WorkspaceID, page, pageSize)
	if err != nil {
		return nil, err
	}

	return &domain.PromptList{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (s *PromptService) UpdatePrompt(ctx context.Context, actor *domain.Actor, id int64, req *UpdatePromptReq) (*domain.Prompt, error) {
	p, err := s.GetPrompt(ctx, actor, id)
	if err != nil {
		return nil, err
	}

	if !p.CanEdit(actor) {
		return nil, ErrNoPermission
	}

	if req.Title != nil {
		p.Title = *req.Title
	}
	if req.Slug != nil {
		p.Slug = *req.Slug
	}
	if req.Description != nil {
		p.Description = *req.Description
	}
	if req.Body != nil {
		p.Body = *req.Body
	}
	if req.MessageFormat != nil {
		p.MessageFormat = *req.MessageFormat
	}
	if req.Visibility != nil {
		p.Visibility = *req.Visibility
	}
	if req.Status != nil {
		p.Status = *req.Status
	}
	if req.TargetProvider != nil {
		p.TargetProvider = *req.TargetProvider
	}
	if req.TargetModel != nil {
		p.TargetModel = *req.TargetModel
	}
	if req.DefaultTemperature != nil {
		p.DefaultTemperature = req.DefaultTemperature
	}
	if req.DefaultMaxTokens != nil {
		p.DefaultMaxTokens = req.DefaultMaxTokens
	}
	if req.UsageNotes != nil {
		p.UsageNotes = *req.UsageNotes
	}

	if err := s.promptRepo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *PromptService) ArchivePrompt(ctx context.Context, actor *domain.Actor, id int64) error {
	p, err := s.GetPrompt(ctx, actor, id)
	if err != nil {
		return err
	}
	if !p.CanEdit(actor) {
		return ErrNoPermission
	}
	p.Status = "archived"
	return s.promptRepo.Update(ctx, p)
}

func (s *PromptService) RestorePrompt(ctx context.Context, actor *domain.Actor, id int64) error {
	p, err := s.GetPrompt(ctx, actor, id)
	if err != nil {
		return err
	}
	if !p.CanEdit(actor) {
		return ErrNoPermission
	}
	p.Status = "draft"
	return s.promptRepo.Update(ctx, p)
}

func generateSlug(title string) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%s-%d", title, r.Intn(10000))
}
