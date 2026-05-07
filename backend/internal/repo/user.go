package repo

import (
	"context"

	"github.com/mymikasa/prompthub/internal/domain"
)

type UserRepo interface {
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, id int64) (*domain.User, error)
}
