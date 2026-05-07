package repo

import "github.com/mymikasa/prompthub/internal/domain"

type WorkspaceRepo interface {
	Create(ws *domain.Workspace) error
	AddMember(member *domain.WorkspaceMember) error
	FindByUserID(userID int64) (*domain.Workspace, error)
}
