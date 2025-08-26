package repository

import (
	"github.com/miffyG/golearn/task4/models"
	"gorm.io/gorm"
)

type UserRepo struct{ db *gorm.DB }

func NewUserRepository(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepo) GetByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("user_name = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
