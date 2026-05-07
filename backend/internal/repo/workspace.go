package repo

import (
	"context"

	"github.com/mymikasa/prompthub/internal/domain"
	"github.com/mymikasa/prompthub/internal/repo/dao"
	"github.com/mymikasa/prompthub/internal/repo/dao/model"
)

type WorkspaceRepo interface {
	Create(ctx context.Context, ws *domain.Workspace) error
	AddMember(ctx context.Context, member *domain.WorkspaceMember) error
	FindByUserID(ctx context.Context, userID int64) (*domain.Workspace, error)
}

type workspaceRepo struct {
	dao dao.WorkspaceDAO
}

func NewWorkspaceRepo(d dao.WorkspaceDAO) WorkspaceRepo {
	return &workspaceRepo{dao: d}
}

func (r *workspaceRepo) Create(ctx context.Context, ws *domain.Workspace) error {
	return r.dao.Create(ctx, toModelWorkspace(ws))
}

func (r *workspaceRepo) AddMember(ctx context.Context, member *domain.WorkspaceMember) error {
	return r.dao.AddMember(ctx, &model.WorkspaceMember{
		WorkspaceID: member.WorkspaceID,
		UserID:      member.UserID,
		Role:        member.Role,
	})
}

func (r *workspaceRepo) FindByUserID(ctx context.Context, userID int64) (*domain.Workspace, error) {
	member, err := r.dao.FindMemberByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	m, err := r.dao.FindByID(ctx, member.WorkspaceID)
	if err != nil {
		return nil, err
	}
	return toDomainWorkspace(m), nil
}

func toModelWorkspace(ws *domain.Workspace) *model.Workspace {
	return &model.Workspace{
		ID:        ws.ID,
		Name:      ws.Name,
		OwnerID:   ws.OwnerID,
		CreatedAt: ws.CreatedAt,
		UpdatedAt: ws.UpdatedAt,
	}
}

func toDomainWorkspace(m *model.Workspace) *domain.Workspace {
	return &domain.Workspace{
		ID:        m.ID,
		Name:      m.Name,
		OwnerID:   m.OwnerID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
