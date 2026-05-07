package dao

import (
	"context"

	"github.com/mymikasa/prompthub/internal/domain"
	"github.com/mymikasa/prompthub/internal/repo"
	"gorm.io/gorm"
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) repo.UserRepo {
	return &UserDAO{db: db}
}

func (d *UserDAO) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	if err := d.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *UserDAO) Create(ctx context.Context, user *domain.User) error {
	return d.db.WithContext(ctx).Create(user).Error
}

func (d *UserDAO) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	var user domain.User
	if err := d.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
