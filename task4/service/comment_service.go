package service

import (
	"github.com/miffyG/golearn/task4/models"
	"github.com/miffyG/golearn/task4/repository"
)

type CommentService struct {
	repo *repository.CommentRepository
}

func NewCommentService(repo *repository.CommentRepository) *CommentService {
	return &CommentService{repo: repo}
}

func (s *CommentService) Create(comment *models.Comment) error {
	return s.repo.Create(comment)
}

func (s *CommentService) GetByPostId(postId uint) ([]models.Comment, error) {
	return s.repo.GetByPostId(postId)
}
