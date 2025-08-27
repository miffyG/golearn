package service

import (
	"time"

	"github.com/miffyG/golearn/task4/internal/models/entity"
	"github.com/miffyG/golearn/task4/internal/repository"
	"github.com/miffyG/golearn/task4/internal/utils"
	"github.com/miffyG/golearn/task4/pkg/config"
)

type UserService struct {
	repo *repository.UserRepo
}

func NewUserService(r *repository.UserRepo) *UserService {
	return &UserService{
		repo: r,
	}
}

func (s *UserService) Register(user *entity.User) error {
	if err := user.SetPassword(user.Password); err != nil {
		return err
	}
	return s.repo.Create(user)
}

func (s *UserService) Login(username, password string) (string, *entity.User, error) {
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		return "", nil, err
	}
	if user == nil {
		return "", nil, nil
	}
	if err := user.CheckPassword(password); err != nil {
		return "", nil, err
	}
	token, err := utils.GenerateJWTToken(config.GetSecretConfig().JwtSecret, user.ID, user.UserName, time.Hour*24)
	if err != nil {
		return "", nil, err
	}
	return token, user, nil
}
