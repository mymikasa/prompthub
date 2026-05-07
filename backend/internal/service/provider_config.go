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
	Name         string `json:"name" binding:"required,min=1,max=100"`
	ProviderType string `json:"providerType" binding:"required"`
	BaseURL      string `json:"baseUrl" binding:"required"`
	APIKey       string `json:"apiKey"`
	DefaultModel string `json:"defaultModel"`
}

type ProviderConfigResponse struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	ProviderType string `json:"providerType"`
	BaseURL      string `json:"baseUrl"`
	HasAPIKey    bool   `json:"hasApiKey"`
	DefaultModel string `json:"defaultModel"`
}

func (s *ProviderConfigService) List(ctx context.Context, actor *domain.Actor) ([]ProviderConfigResponse, error) {
	list, err := s.repo.ListByWorkspaceID(ctx, actor.WorkspaceID)
	if err != nil {
		return nil, err
	}
	result := make([]ProviderConfigResponse, len(list))
	for i, pc := range list {
		result[i] = ProviderConfigResponse{
			ID:           pc.ID,
			Name:         pc.Name,
			ProviderType: pc.ProviderType,
			BaseURL:      pc.BaseURL,
			HasAPIKey:    pc.EncryptedAPIKey != "",
			DefaultModel: pc.DefaultModel,
		}
	}
	return result, nil
}

func (s *ProviderConfigService) GetByID(ctx context.Context, actor *domain.Actor, id int64) (*ProviderConfigResponse, error) {
	pc, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProviderConfigNotFound
		}
		return nil, err
	}
	if pc.WorkspaceID != actor.WorkspaceID {
		return nil, ErrProviderConfigNotFound
	}
	return &ProviderConfigResponse{
		ID:           pc.ID,
		Name:         pc.Name,
		ProviderType: pc.ProviderType,
		BaseURL:      pc.BaseURL,
		HasAPIKey:    pc.EncryptedAPIKey != "",
		DefaultModel: pc.DefaultModel,
	}, nil
}

func (s *ProviderConfigService) Save(ctx context.Context, actor *domain.Actor, req *SaveProviderConfigReq) (*ProviderConfigResponse, error) {
	existing, err := s.repo.FindByWorkspaceAndType(ctx, actor.WorkspaceID, req.ProviderType)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	pc := &domain.ProviderConfig{
		WorkspaceID:  actor.WorkspaceID,
		ProviderType: req.ProviderType,
		Name:         req.Name,
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

	if existing != nil {
		if err := s.repo.Update(ctx, pc); err != nil {
			return nil, err
		}
	} else {
		if err := s.repo.Create(ctx, pc); err != nil {
			return nil, err
		}
	}

	return &ProviderConfigResponse{
		ID:           pc.ID,
		Name:         pc.Name,
		ProviderType: pc.ProviderType,
		BaseURL:      pc.BaseURL,
		HasAPIKey:    pc.EncryptedAPIKey != "",
		DefaultModel: pc.DefaultModel,
	}, nil
}

func (s *ProviderConfigService) Delete(ctx context.Context, actor *domain.Actor, id int64) error {
	pc, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProviderConfigNotFound
		}
		return err
	}
	if pc.WorkspaceID != actor.WorkspaceID {
		return ErrProviderConfigNotFound
	}
	return s.repo.Delete(ctx, id)
}

func (s *ProviderConfigService) DecryptAPIKey(ctx context.Context, workspaceID int64, providerType string) (string, error) {
	pc, err := s.repo.FindByWorkspaceAndType(ctx, workspaceID, providerType)
	if err != nil {
		return "", ErrProviderConfigNotFound
	}
	if pc.EncryptedAPIKey == "" {
		return "", ErrProviderConfigNotFound
	}
	return crypto.Decrypt(pc.EncryptedAPIKey, s.encKey)
}

func (s *ProviderConfigService) GetProviderConfig(ctx context.Context, workspaceID int64, providerType string) (*domain.ProviderConfig, error) {
	pc, err := s.repo.FindByWorkspaceAndType(ctx, workspaceID, providerType)
	if err != nil {
		return nil, ErrProviderConfigNotFound
	}
	return pc, nil
}

func (s *ProviderConfigService) GetProviderConfigByID(ctx context.Context, actor *domain.Actor, id int64) (*domain.ProviderConfig, error) {
	pc, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrProviderConfigNotFound
	}
	if pc.WorkspaceID != actor.WorkspaceID {
		return nil, ErrProviderConfigNotFound
	}
	return pc, nil
}

func (s *ProviderConfigService) DecryptAPIKeyByConfig(ctx context.Context, pc *domain.ProviderConfig) (string, error) {
	if pc.EncryptedAPIKey == "" {
		return "", ErrProviderConfigNotFound
	}
	return crypto.Decrypt(pc.EncryptedAPIKey, s.encKey)
}
