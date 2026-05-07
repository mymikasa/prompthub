package dao

import (
	"context"

	"github.com/mymikasa/prompthub/internal/repo/dao/model"
	"gorm.io/gorm"
)

type UserDAO interface {
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	Create(ctx context.Context, m *model.User) error
	FindByID(ctx context.Context, id int64) (*model.User, error)
}

type userDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) UserDAO {
	return &userDAO{db: db}
}

func (d *userDAO) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var m model.User
	if err := d.db.WithContext(ctx).Where("email = ?", email).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (d *userDAO) Create(ctx context.Context, m *model.User) error {
	return d.db.WithContext(ctx).Create(m).Error
}

func (d *userDAO) FindByID(ctx context.Context, id int64) (*model.User, error) {
	var m model.User
	if err := d.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}
