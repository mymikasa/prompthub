package repo

import (
	"context"

	"github.com/mymikasa/prompthub/internal/domain"
)

type WorkspaceRepo interface {
	Create(ctx context.Context, ws *domain.Workspace) error
	AddMember(ctx context.Context, member *domain.WorkspaceMember) error
	FindByUserID(ctx context.Context, userID int64) (*domain.Workspace, error)
}
