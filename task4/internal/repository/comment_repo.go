package repository

import (
	"github.com/miffyG/golearn/task4/internal/models/entity"
	"gorm.io/gorm"
)

type CommentRepository struct{ db *gorm.DB }

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) Create(comment *entity.Comment) error {
	return r.db.Create(comment).Error
}

func (r *CommentRepository) GetByPostId(postId uint) ([]entity.Comment, error) {
	var comments []entity.Comment
	if err := r.db.Where("post_id = ?", postId).Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}
