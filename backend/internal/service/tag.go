package service

import (
	"context"

	"github.com/mymikasa/prompthub/internal/domain"
	"github.com/mymikasa/prompthub/internal/repo"
)

type TagService struct {
	tagRepo repo.TagRepo
}

func NewTagService(tagRepo repo.TagRepo) *TagService {
	return &TagService{tagRepo: tagRepo}
}

func (s *TagService) ListTags(ctx context.Context, actor *domain.Actor) ([]string, error) {
	tags, err := s.tagRepo.FindByWorkspace(ctx, actor.WorkspaceID)
	if err != nil {
		return nil, err
	}
	names := make([]string, len(tags))
	for i, t := range tags {
		names[i] = t.Name
	}
	return names, nil
}
