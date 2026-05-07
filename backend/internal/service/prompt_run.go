package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/mymikasa/prompthub/internal/domain"
	"github.com/mymikasa/prompthub/internal/repo"
	oai "github.com/mymikasa/prompthub/pkg/openai"
)

var (
	ErrNoProviderConfig = errors.New("provider config not found, please configure in settings")
	ErrRunFailed        = errors.New("prompt run failed")
)

type RunPromptReq struct {
	Variables   map[string]string `json:"variables"`
	TestCaseID  *int64            `json:"testCaseId"`
	ProviderID  *int64            `json:"providerId"`
	Model       *string           `json:"model"`
}

type RunPromptResult struct {
	Run *domain.PromptRun `json:"run"`
}

type PromptRunService struct {
	runRepo        repo.PromptRunRepo
	promptRepo     repo.PromptRepo
	testCaseRepo   repo.TestCaseRepo
	variableRepo   repo.VariableRepo
	providerSvc    *ProviderConfigService
	promptSvc      *PromptService
	openaiClient   *oai.Client
}

func NewPromptRunService(
	runRepo repo.PromptRunRepo,
	promptRepo repo.PromptRepo,
	testCaseRepo repo.TestCaseRepo,
	variableRepo repo.VariableRepo,
	providerSvc *ProviderConfigService,
	promptSvc *PromptService,
) *PromptRunService {
	return &PromptRunService{
		runRepo:      runRepo,
		promptRepo:   promptRepo,
		testCaseRepo: testCaseRepo,
		variableRepo: variableRepo,
		providerSvc:  providerSvc,
		promptSvc:    promptSvc,
		openaiClient: oai.NewClient(),
	}
}

func (s *PromptRunService) Run(ctx context.Context, actor *domain.Actor, promptID int64, req *RunPromptReq) (*RunPromptResult, error) {
	p, err := s.promptSvc.GetPrompt(ctx, actor, promptID)
	if err != nil {
		return nil, err
	}
	if !p.CanEdit(actor) {
		return nil, ErrNoPermission
	}

	variables := req.Variables
	if req.TestCaseID != nil {
		tc, err := s.testCaseRepo.FindByID(ctx, *req.TestCaseID)
		if err != nil {
			return nil, ErrTestCaseNotFound
		}
		if tc.PromptID != promptID {
			return nil, ErrTestCaseNotFound
		}
		variables = tc.VariableValues
	}

	if variables == nil {
		variables = map[string]string{}
	}

	renderReq := &RenderReq{Variables: variables}
	rendered, err := s.promptSvc.RenderPrompt(ctx, actor, promptID, renderReq)
	if err != nil {
		return nil, err
	}

	var providerCfg *domain.ProviderConfig
	if req.ProviderID != nil {
		providerCfg, err = s.providerSvc.GetProviderConfigByID(ctx, actor, *req.ProviderID)
	} else {
		providerCfg, err = s.providerSvc.GetProviderConfig(ctx, actor.WorkspaceID, "openai_compatible")
	}
	if err != nil {
		return nil, ErrNoProviderConfig
	}

	apiKey, err := s.providerSvc.DecryptAPIKeyByConfig(ctx, providerCfg)
	if err != nil {
		return nil, ErrNoProviderConfig
	}

	model := providerCfg.DefaultModel
	if req.Model != nil && *req.Model != "" {
		model = *req.Model
	}
	if model == "" {
		return nil, errors.New("no model specified, configure default model in provider settings")
	}

	var messages []oai.ChatMessage
	if rendered.Messages != nil {
		for _, m := range rendered.Messages {
			messages = append(messages, oai.ChatMessage{Role: m.Role, Content: m.Content})
		}
	} else {
		messages = []oai.ChatMessage{{Role: "user", Content: rendered.Content}}
	}

	chatReq := &oai.ChatRequest{Messages: messages}
	result, err := s.openaiClient.Chat(ctx, providerCfg.BaseURL, apiKey, model, chatReq, p.DefaultTemperature, p.DefaultMaxTokens)

	run := &domain.PromptRun{
		PromptID:               promptID,
		PromptVersionID:        p.CurrentVersionID,
		TestCaseID:             req.TestCaseID,
		Provider:               providerCfg.ProviderType,
		Model:                  model,
		InputVariables:         variables,
		RenderedPromptSnapshot: rendered,
		CreatedBy:              actor.UserID,
	}

	if err != nil {
		slog.Error("prompt run failed",
			slog.Int64("prompt_id", promptID),
			slog.String("model", model),
			slog.String("error", err.Error()),
		)
		run.ErrorMessage = "model request failed, please check provider configuration"
		run.Latency = 0
	} else {
		run.OutputText = result.Content
		run.Latency = result.LatencyMs
		run.TokenUsage = &domain.TokenUsage{
			PromptTokens:     result.Usage.PromptTokens,
			CompletionTokens: result.Usage.CompletionTokens,
			TotalTokens:      result.Usage.TotalTokens,
		}
	}

	if saveErr := s.runRepo.Create(ctx, run); saveErr != nil {
		slog.Error("save run record failed", slog.String("error", saveErr.Error()))
		if err != nil {
			return nil, ErrRunFailed
		}
		return nil, saveErr
	}

	if err != nil {
		return nil, ErrRunFailed
	}

	return &RunPromptResult{Run: run}, nil
}

func (s *PromptRunService) ListRuns(ctx context.Context, actor *domain.Actor, promptID int64, page, pageSize int) ([]*domain.PromptRun, int64, error) {
	if _, err := s.promptSvc.GetPrompt(ctx, actor, promptID); err != nil {
		return nil, 0, err
	}
	return s.runRepo.FindByPromptID(ctx, promptID, page, pageSize)
}
