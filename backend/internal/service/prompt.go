package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/mymikasa/prompthub/internal/domain"
	"github.com/mymikasa/prompthub/internal/repo"
	"github.com/mymikasa/prompthub/pkg/vars"
	"gorm.io/gorm"
)

var (
	ErrPromptNotFound  = errors.New("prompt not found")
	ErrNoPermission    = errors.New("no permission")
	ErrMissingRequired = errors.New("missing required variable")
)

type PromptService struct {
	promptRepo       repo.PromptRepo
	tagRepo          repo.TagRepo
	variableRepo     repo.VariableRepo
	versionRepo      repo.VersionRepo
	providerConfigRepo repo.ProviderConfigRepo
}

func NewPromptService(promptRepo repo.PromptRepo, tagRepo repo.TagRepo, variableRepo repo.VariableRepo, versionRepo repo.VersionRepo, providerConfigRepo repo.ProviderConfigRepo) *PromptService {
	return &PromptService{promptRepo: promptRepo, tagRepo: tagRepo, variableRepo: variableRepo, versionRepo: versionRepo, providerConfigRepo: providerConfigRepo}
}

type PromptFilter struct {
	Keyword  string
	Statuses []string
	Tags     []string
	Provider string
	Model    string
}

type CreatePromptReq struct {
	Title              string   `json:"title" binding:"required,min=1,max=200"`
	Slug               string   `json:"slug"`
	Description        string   `json:"description"`
	Body               string   `json:"body" binding:"required"`
	MessageFormat      string   `json:"messageFormat" binding:"required,oneof=single_text chat_messages"`
	Visibility         string   `json:"visibility" binding:"required,oneof=private workspace"`
	ProviderConfigID   *int64   `json:"providerConfigId"`
	TargetProvider     string   `json:"targetProvider"`
	TargetModel        string   `json:"targetModel"`
	DefaultTemperature *float64 `json:"defaultTemperature"`
	DefaultMaxTokens   *int     `json:"defaultMaxTokens"`
	UsageNotes         string   `json:"usageNotes"`
	Tags               []string `json:"tags"`
}

type UpdatePromptReq struct {
	Title              *string  `json:"title"`
	Slug               *string  `json:"slug"`
	Description        *string  `json:"description"`
	Body               *string  `json:"body"`
	MessageFormat      *string  `json:"messageFormat"`
	Visibility         *string  `json:"visibility"`
	Status             *string  `json:"status"`
	ProviderConfigID   **int64  `json:"providerConfigId"`
	TargetProvider     *string  `json:"targetProvider"`
	TargetModel        *string  `json:"targetModel"`
	DefaultTemperature *float64 `json:"defaultTemperature"`
	DefaultMaxTokens   *int     `json:"defaultMaxTokens"`
	UsageNotes         *string  `json:"usageNotes"`
	Tags               []string `json:"tags"`
}

func (s *PromptService) CreatePrompt(ctx context.Context, actor *domain.Actor, req *CreatePromptReq) (*domain.Prompt, error) {
	slug := req.Slug
	if slug == "" {
		slug = generateSlug(req.Title)
	}

	targetProvider := req.TargetProvider
	targetModel := req.TargetModel
	var providerConfigID *int64

	if req.ProviderConfigID != nil {
		pc, err := s.providerConfigRepo.FindByID(ctx, *req.ProviderConfigID)
		if err != nil {
			return nil, fmt.Errorf("provider config not found")
		}
		if pc.WorkspaceID != actor.WorkspaceID {
			return nil, fmt.Errorf("provider config not found")
		}
		providerConfigID = req.ProviderConfigID
		targetProvider = pc.ProviderType
		if targetModel == "" {
			targetModel = pc.DefaultModel
		}
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
		TargetProvider:     targetProvider,
		TargetModel:        targetModel,
		ProviderConfigID:   providerConfigID,
		DefaultTemperature: req.DefaultTemperature,
		DefaultMaxTokens:   req.DefaultMaxTokens,
		UsageNotes:         req.UsageNotes,
		Tags:               req.Tags,
	}

	if err := s.promptRepo.Create(ctx, p); err != nil {
		return nil, err
	}

	if err := s.syncTags(ctx, actor.WorkspaceID, p.ID, req.Tags); err != nil {
		return nil, err
	}

	if err := s.syncVariables(ctx, p.ID, req.Body); err != nil {
		return nil, err
	}

	if err := s.createInitialVersion(ctx, p, req.Tags); err != nil {
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

	p.Tags, err = s.loadTagNames(ctx, p.ID)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (s *PromptService) ListPrompts(ctx context.Context, actor *domain.Actor, page, pageSize int, filter PromptFilter) (*domain.PromptList, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	repoFilter := repo.PromptFilter{
		Keyword:  filter.Keyword,
		Statuses: filter.Statuses,
		TagNames: filter.Tags,
		Provider: filter.Provider,
		Model:    filter.Model,
	}

	items, total, err := s.promptRepo.FindByWorkspace(ctx, actor.WorkspaceID, page, pageSize, repoFilter)
	if err != nil {
		return nil, err
	}

	for _, p := range items {
		p.Tags, err = s.loadTagNames(ctx, p.ID)
		if err != nil {
			return nil, err
		}
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
	if req.ProviderConfigID != nil {
		configID := *req.ProviderConfigID
		if configID != nil {
			pc, err := s.providerConfigRepo.FindByID(ctx, *configID)
			if err != nil {
				return nil, fmt.Errorf("provider config not found")
			}
			if pc.WorkspaceID != actor.WorkspaceID {
				return nil, fmt.Errorf("provider config not found")
			}
			p.ProviderConfigID = configID
			p.TargetProvider = pc.ProviderType
			if req.TargetModel == nil {
				p.TargetModel = pc.DefaultModel
			}
		} else {
			p.ProviderConfigID = nil
		}
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
	if req.Tags != nil {
		p.Tags = req.Tags
	}

	if err := s.promptRepo.Update(ctx, p); err != nil {
		return nil, err
	}

	if req.Tags != nil {
		if err := s.syncTags(ctx, actor.WorkspaceID, p.ID, req.Tags); err != nil {
			return nil, err
		}
	}

	if req.Body != nil {
		if err := s.syncVariables(ctx, p.ID, p.Body); err != nil {
			return nil, err
		}
	}

	if err := s.createUpdateVersion(ctx, p); err != nil {
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

func (s *PromptService) syncTags(ctx context.Context, workspaceID, promptID int64, tagNames []string) error {
	if len(tagNames) == 0 {
		return s.tagRepo.ReplacePromptTags(ctx, promptID, nil)
	}

	tags, err := s.tagRepo.FindOrCreateByNames(ctx, workspaceID, tagNames)
	if err != nil {
		return err
	}

	tagIDs := make([]int64, len(tags))
	for i, t := range tags {
		tagIDs[i] = t.ID
	}
	return s.tagRepo.ReplacePromptTags(ctx, promptID, tagIDs)
}

func (s *PromptService) syncVariables(ctx context.Context, promptID int64, body string) error {
	extracted := vars.Extract(body)

	existing, err := s.variableRepo.FindByPromptID(ctx, promptID)
	if err != nil {
		return err
	}

	existingMap := make(map[string]*domain.PromptVariable, len(existing))
	for _, v := range existing {
		existingMap[v.Name] = v
	}

	var newVars []*domain.PromptVariable
	for _, name := range extracted {
		if _, ok := existingMap[name]; !ok {
			newVars = append(newVars, &domain.PromptVariable{
				PromptID: promptID,
				Name:     name,
			})
		}
	}

	if len(newVars) > 0 {
		if err := s.variableRepo.BatchCreate(ctx, newVars); err != nil {
			return err
		}
	}

	if err := s.variableRepo.DeleteNotIn(ctx, promptID, extracted); err != nil {
		return err
	}

	return nil
}

func (s *PromptService) loadTagNames(ctx context.Context, promptID int64) ([]string, error) {
	tags, err := s.tagRepo.FindByPromptID(ctx, promptID)
	if err != nil {
		return nil, err
	}
	names := make([]string, len(tags))
	for i, t := range tags {
		names[i] = t.Name
	}
	return names, nil
}

type RenderResult struct {
	Content   string        `json:"content"`
	Messages  []ChatMessage `json:"messages"`
	Variables []string      `json:"variables"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type RenderReq struct {
	Variables map[string]string `json:"variables"`
}

func (s *PromptService) RenderPrompt(ctx context.Context, actor *domain.Actor, promptID int64, req *RenderReq) (*RenderResult, error) {
	p, err := s.GetPrompt(ctx, actor, promptID)
	if err != nil {
		return nil, err
	}

	allVars, err := s.variableRepo.FindByPromptID(ctx, promptID)
	if err != nil {
		return nil, err
	}
	var missing []string
	for _, v := range allVars {
		if v.Required {
			val, ok := req.Variables[v.Name]
			if !ok || val == "" {
				if v.DefaultValue == "" {
					missing = append(missing, v.Name)
				}
			}
		}
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("%w: %v", ErrMissingRequired, missing)
	}

	resolved := make(map[string]string)
	for _, v := range allVars {
		if val, ok := req.Variables[v.Name]; ok && val != "" {
			resolved[v.Name] = val
		} else if v.DefaultValue != "" {
			resolved[v.Name] = v.DefaultValue
		}
	}

	varNames := make([]string, len(allVars))
	for i, v := range allVars {
		varNames[i] = v.Name
	}

	result := &RenderResult{Variables: varNames}

	switch p.MessageFormat {
	case "single_text":
		result.Content = replaceVars(p.Body, resolved)
		result.Messages = nil
	case "chat_messages":
		var msgs []ChatMessage
		if err := json.Unmarshal([]byte(p.Body), &msgs); err != nil {
			return nil, fmt.Errorf("invalid chat messages format: %w", err)
		}
		var sb strings.Builder
		for i, m := range msgs {
			m.Content = replaceVars(m.Content, resolved)
			msgs[i] = m
			sb.WriteString(m.Content)
			if i < len(msgs)-1 {
				sb.WriteString("\n")
			}
		}
		result.Messages = msgs
		result.Content = sb.String()
	}

	return result, nil
}

func replaceVars(text string, vars map[string]string) string {
	for name, val := range vars {
		text = strings.ReplaceAll(text, "{{"+name+"}}", val)
	}
	return text
}

func generateSlug(title string) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%s-%d", title, r.Intn(10000))
}
