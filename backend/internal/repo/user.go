package repo

import (
	"context"

	"github.com/mymikasa/prompthub/internal/domain"
	"github.com/mymikasa/prompthub/internal/repo/dao"
	"github.com/mymikasa/prompthub/internal/repo/dao/model"
)

type UserRepo interface {
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, id int64) (*domain.User, error)
}

type userRepo struct {
	dao dao.UserDAO
}

func NewUserRepo(d dao.UserDAO) UserRepo {
	return &userRepo{dao: d}
}

func (r *userRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	m, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return toDomainUser(m), nil
}

func (r *userRepo) Create(ctx context.Context, user *domain.User) error {
	return r.dao.Create(ctx, toModelUser(user))
}

func (r *userRepo) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	m, err := r.dao.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toDomainUser(m), nil
}

func toModelUser(u *domain.User) *model.User {
	return &model.User{
		ID:           u.ID,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Nickname:     u.Nickname,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

func toDomainUser(m *model.User) *domain.User {
	return &domain.User{
		ID:           m.ID,
		Email:        m.Email,
		PasswordHash: m.PasswordHash,
		Nickname:     m.Nickname,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}
