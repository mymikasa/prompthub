package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/mymikasa/prompthub/internal/domain"
	"github.com/mymikasa/prompthub/internal/repo"
	"github.com/mymikasa/prompthub/pkg/crypto"
	"gorm.io/gorm"
)

var ErrProviderConfigNotFound = errors.New("provider config not found")

type ProviderConfigService struct {
	repo   repo.ProviderConfigRepo
	encKey []byte
}

func NewProviderConfigService(repo repo.ProviderConfigRepo, encryptionKeyHex string) (*ProviderConfigService, error) {
	key, err := crypto.DeriveKey(encryptionKeyHex)
	if err != nil {
		return nil, err
	}
	return &ProviderConfigService{repo: repo, encKey: key}, nil
}

type SaveProviderConfigReq struct {
	ProviderType string `json:"providerType" binding:"required,oneof=openai_compatible"`
	BaseURL      string `json:"baseUrl" binding:"required"`
	APIKey       string `json:"apiKey"`
	DefaultModel string `json:"defaultModel"`
}

type ProviderConfigResponse struct {
	ID           int64  `json:"id"`
	ProviderType string `json:"providerType"`
	BaseURL      string `json:"baseUrl"`
	HasAPIKey    bool   `json:"hasApiKey"`
	DefaultModel string `json:"defaultModel"`
}

func (s *ProviderConfigService) Get(ctx context.Context, actor *domain.Actor) (*ProviderConfigResponse, error) {
	pc, err := s.repo.FindByWorkspaceID(ctx, actor.WorkspaceID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &ProviderConfigResponse{}, nil
		}
		return nil, err
	}
	return &ProviderConfigResponse{
		ID:           pc.ID,
		ProviderType: pc.ProviderType,
		BaseURL:      pc.BaseURL,
		HasAPIKey:    pc.HasAPIKey,
		DefaultModel: pc.DefaultModel,
	}, nil
}

func (s *ProviderConfigService) Save(ctx context.Context, actor *domain.Actor, req *SaveProviderConfigReq) (*ProviderConfigResponse, error) {
	existing, err := s.repo.FindByWorkspaceID(ctx, actor.WorkspaceID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	pc := &domain.ProviderConfig{
		WorkspaceID:  actor.WorkspaceID,
		ProviderType: req.ProviderType,
		BaseURL:      req.BaseURL,
		DefaultModel: req.DefaultModel,
		CreatedBy:    actor.UserID,
	}

	if existing != nil {
		pc.ID = existing.ID
		pc.CreatedBy = existing.CreatedBy
	}

	if req.APIKey != "" {
		encrypted, err := crypto.Encrypt(req.APIKey, s.encKey)
		if err != nil {
			slog.Error("encrypt api key failed", slog.String("error", err.Error()))
			return nil, errors.New("failed to save api key")
		}
		pc.EncryptedAPIKey = encrypted
	} else if existing != nil {
		pc.EncryptedAPIKey = existing.EncryptedAPIKey
	}

	if err := s.repo.Save(ctx, pc); err != nil {
		return nil, err
	}

	return &ProviderConfigResponse{
		ID:           pc.ID,
		ProviderType: pc.ProviderType,
		BaseURL:      pc.BaseURL,
		HasAPIKey:    pc.EncryptedAPIKey != "",
		DefaultModel: pc.DefaultModel,
	}, nil
}

func (s *ProviderConfigService) DecryptAPIKey(ctx context.Context, workspaceID int64) (string, error) {
	pc, err := s.repo.FindByWorkspaceID(ctx, workspaceID)
	if err != nil {
		return "", ErrProviderConfigNotFound
	}
	if pc.EncryptedAPIKey == "" {
		return "", ErrProviderConfigNotFound
	}
	return crypto.Decrypt(pc.EncryptedAPIKey, s.encKey)
}
