package service

import (
	"context"
	"errors"

	"github.com/mymikasa/prompthub/internal/domain"
	"github.com/mymikasa/prompthub/internal/repo"
	"gorm.io/gorm"
)

var ErrVersionNotFound = errors.New("version not found")

type VersionService struct {
	versionRepo repo.VersionRepo
	promptRepo  repo.PromptRepo
	tagRepo     repo.TagRepo
	variableRepo repo.VariableRepo
}

func NewVersionService(versionRepo repo.VersionRepo, promptRepo repo.PromptRepo, tagRepo repo.TagRepo, variableRepo repo.VariableRepo) *VersionService {
	return &VersionService{versionRepo: versionRepo, promptRepo: promptRepo, tagRepo: tagRepo, variableRepo: variableRepo}
}

func (s *VersionService) ListVersions(ctx context.Context, actor *domain.Actor, promptID int64) ([]*domain.PromptVersion, error) {
	p, err := s.promptRepo.FindByID(ctx, promptID)
	if err != nil {
		return nil, ErrPromptNotFound
	}
	if p.WorkspaceID != actor.WorkspaceID || !p.CanView(actor.UserID) {
		return nil, ErrPromptNotFound
	}
	return s.versionRepo.FindByPromptID(ctx, promptID)
}

func (s *VersionService) GetVersion(ctx context.Context, actor *domain.Actor, promptID, versionID int64) (*domain.PromptVersion, error) {
	p, err := s.promptRepo.FindByID(ctx, promptID)
	if err != nil {
		return nil, ErrPromptNotFound
	}
	if p.WorkspaceID != actor.WorkspaceID || !p.CanView(actor.UserID) {
		return nil, ErrPromptNotFound
	}

	v, err := s.versionRepo.FindByID(ctx, versionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrVersionNotFound
		}
		return nil, err
	}
	if v.PromptID != promptID {
		return nil, ErrVersionNotFound
	}
	return v, nil
}

func (s *VersionService) RestoreVersion(ctx context.Context, actor *domain.Actor, promptID, versionID int64) (*domain.PromptVersion, error) {
	p, err := s.promptRepo.FindByID(ctx, promptID)
	if err != nil {
		return nil, ErrPromptNotFound
	}
	if p.WorkspaceID != actor.WorkspaceID || !p.CanEdit(actor) {
		return nil, ErrNoPermission
	}

	v, err := s.versionRepo.FindByID(ctx, versionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrVersionNotFound
		}
		return nil, err
	}
	if v.PromptID != promptID {
		return nil, ErrVersionNotFound
	}

	// restore snapshot into prompt
	snap := v.Snapshot
	p.Body = snap.Content
	p.TargetProvider = snap.TargetProvider
	p.TargetModel = snap.TargetModel
	p.Status = snap.Status
	if err := s.promptRepo.Update(ctx, p); err != nil {
		return nil, err
	}

	// sync tags
	if err := s.syncTagsFromNames(ctx, actor.WorkspaceID, p.ID, snap.Tags); err != nil {
		return nil, err
	}

	// create new version from restored state
	newVer, err := s.createVersion(ctx, p, snap.Tags, "restored from version "+string(rune('0'+v.VersionNumber)))
	if err != nil {
		return nil, err
	}

	return newVer, nil
}

func (s *VersionService) createVersion(ctx context.Context, p *domain.Prompt, tagNames []string, changeNote string) (*domain.PromptVersion, error) {
	latest, err := s.versionRepo.LatestVersion(ctx, p.ID)
	if err != nil {
		return nil, err
	}

	vars, err := s.variableRepo.FindByPromptID(ctx, p.ID)
	if err != nil {
		return nil, err
	}

	snapshot := domain.VersionSnapshot{
		Content:        p.Body,
		TargetProvider: p.TargetProvider,
		TargetModel:    p.TargetModel,
		Status:         p.Status,
		Tags:           tagNames,
	}
	if len(vars) > 0 {
		for _, v := range vars {
			snapshot.Variables = append(snapshot.Variables, *v)
		}
	}

	v := &domain.PromptVersion{
		PromptID:      p.ID,
		VersionNumber: latest + 1,
		ChangeNote:    changeNote,
		Snapshot:      snapshot,
	}

	if err := s.versionRepo.Create(ctx, v); err != nil {
		return nil, err
	}

	p.CurrentVersionID = v.ID
	if err := s.promptRepo.Update(ctx, p); err != nil {
		return nil, err
	}

	return v, nil
}

func (s *VersionService) syncTagsFromNames(ctx context.Context, workspaceID, promptID int64, names []string) error {
	if len(names) == 0 {
		return s.tagRepo.ReplacePromptTags(ctx, promptID, nil)
	}
	tags, err := s.tagRepo.FindOrCreateByNames(ctx, workspaceID, names)
	if err != nil {
		return err
	}
	tagIDs := make([]int64, len(tags))
	for i, t := range tags {
		tagIDs[i] = t.ID
	}
	return s.tagRepo.ReplacePromptTags(ctx, promptID, tagIDs)
}
