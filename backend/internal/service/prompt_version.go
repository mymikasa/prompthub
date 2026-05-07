package service

import (
	"context"

	"github.com/mymikasa/prompthub/internal/domain"
)

func (s *PromptService) createInitialVersion(ctx context.Context, p *domain.Prompt, tagNames []string) error {
	variables, err := s.variableRepo.FindByPromptID(ctx, p.ID)
	if err != nil {
		return err
	}

	snap := buildSnapshot(p, variables, tagNames)
	v := &domain.PromptVersion{
		PromptID:      p.ID,
		VersionNumber: 1,
		Snapshot:      snap,
	}

	if err := s.versionRepo.Create(ctx, v); err != nil {
		return err
	}

	p.CurrentVersionID = v.ID
	return s.promptRepo.Update(ctx, p)
}

func (s *PromptService) createUpdateVersion(ctx context.Context, p *domain.Prompt) error {
	tagNames, err := s.loadTagNames(ctx, p.ID)
	if err != nil {
		return err
	}

	variables, err := s.variableRepo.FindByPromptID(ctx, p.ID)
	if err != nil {
		return err
	}

	latest, err := s.versionRepo.LatestVersion(ctx, p.ID)
	if err != nil {
		return err
	}

	snap := buildSnapshot(p, variables, tagNames)
	v := &domain.PromptVersion{
		PromptID:      p.ID,
		VersionNumber: latest + 1,
		Snapshot:      snap,
	}

	if err := s.versionRepo.Create(ctx, v); err != nil {
		return err
	}

	p.CurrentVersionID = v.ID
	return s.promptRepo.Update(ctx, p)
}

func buildSnapshot(p *domain.Prompt, variables []*domain.PromptVariable, tagNames []string) domain.VersionSnapshot {
	snap := domain.VersionSnapshot{
		Content:        p.Body,
		TargetProvider: p.TargetProvider,
		TargetModel:    p.TargetModel,
		Status:         p.Status,
		Tags:           tagNames,
	}
	for _, v := range variables {
		snap.Variables = append(snap.Variables, *v)
	}
	return snap
}
