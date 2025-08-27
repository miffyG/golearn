package service

import (
	"github.com/miffyG/golearn/task4/internal/models/entity"
	"github.com/miffyG/golearn/task4/internal/repository"
)

type CommentService struct {
	repo *repository.CommentRepository
}

func NewCommentService(repo *repository.CommentRepository) *CommentService {
	return &CommentService{repo: repo}
}

func (s *CommentService) Create(comment *entity.Comment) error {
	return s.repo.Create(comment)
}

func (s *CommentService) GetByPostId(postId uint) ([]entity.Comment, error) {
	return s.repo.GetByPostId(postId)
}
