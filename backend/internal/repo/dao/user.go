package dao

import (
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

func (d *UserDAO) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	if err := d.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *UserDAO) Create(user *domain.User) error {
	return d.db.Create(user).Error
}

func (d *UserDAO) FindByID(id int64) (*domain.User, error) {
	var user domain.User
	if err := d.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
