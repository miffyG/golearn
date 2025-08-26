package service

import (
	"errors"

	"github.com/miffyG/golearn/task4/models"
	"github.com/miffyG/golearn/task4/repository"
)

type PostService struct {
	repo *repository.PostRepository
}

func NewPostService(repo *repository.PostRepository) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) Create(post *models.Post) error {
	return s.repo.Create(post)
}

func (s *PostService) GetAll() ([]models.Post, error) {
	return s.repo.GetAll()
}

func (s *PostService) GetByID(id uint) (*models.Post, error) {
	return s.repo.GetByID(id)
}

func (s *PostService) Update(userId uint, post *models.Post) error {
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
