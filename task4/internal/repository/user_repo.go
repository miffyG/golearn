package repository

import (
	"github.com/miffyG/golearn/task4/internal/models/entity"
	"gorm.io/gorm"
)

type UserRepo struct{ db *gorm.DB }

func NewUserRepository(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(user *entity.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepo) GetByUsername(username string) (*entity.User, error) {
	var user entity.User
	if err := r.db.Where("user_name = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
