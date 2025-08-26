package repository

import (
	"github.com/miffyG/golearn/task4/models"
	"gorm.io/gorm"
)

type CommentRepository struct{ db *gorm.DB }

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) Create(comment *models.Comment) error {
	return r.db.Create(comment).Error
}

func (r *CommentRepository) GetByPostId(postId uint) ([]models.Comment, error) {
	var comments []models.Comment
	if err := r.db.Where("post_id = ?", postId).Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}
