package repo

import "github.com/mymikasa/prompthub/internal/domain"

type UserRepo interface {
	FindByEmail(email string) (*domain.User, error)
	Create(user *domain.User) error
	FindByID(id int64) (*domain.User, error)
}
