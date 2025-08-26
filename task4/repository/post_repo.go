package repository

import (
	"github.com/miffyG/golearn/task4/models"
	"gorm.io/gorm"
)

type PostRepository struct{ db *gorm.DB }

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(post *models.Post) error {
	return r.db.Create(post).Error
}

func (r *PostRepository) GetAll() ([]models.Post, error) {
	var posts []models.Post
	if err := r.db.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "user_name")
	}).Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *PostRepository) GetByID(id uint) (*models.Post, error) {
	var post models.Post
	if err := r.db.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "user_name")
	}).Preload("Comments", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "content", "user_id", "post_id")
	}).First(&post, id).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *PostRepository) Update(post *models.Post) error {
	return r.db.Save(post).Error
}

func (r *PostRepository) Delete(post *models.Post) error {
	return r.db.Delete(post).Error
}
