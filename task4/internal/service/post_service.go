package service

import (
	"errors"

	"github.com/miffyG/golearn/task4/internal/models/entity"
	"github.com/miffyG/golearn/task4/internal/repository"
)

type PostService struct {
	repo *repository.PostRepository
}

func NewPostService(repo *repository.PostRepository) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) Create(post *entity.Post) error {
	return s.repo.Create(post)
}

func (s *PostService) GetAll() ([]entity.Post, error) {
	return s.repo.GetAll()
}

func (s *PostService) GetByID(id uint) (*entity.Post, error) {
	return s.repo.GetByID(id)
}

func (s *PostService) Update(userId uint, post *entity.Post) error {
	p, err := s.repo.GetByID(post.ID)
	if err != nil {
		return err
	}
	if p.UserID != userId {
		return errors.New("unauthorized")
	}
	p.Title = post.Title
	p.Content = post.Content
	return s.repo.Update(p)
}

func (s *PostService) Delete(userId, postId uint) error {
	p, err := s.repo.GetByID(postId)
	if err != nil {
		return err
	}
	if p.UserID != userId {
		return errors.New("unauthorized")
	}
	return s.repo.Delete(p)
}
