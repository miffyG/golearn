package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// users 表：存储用户信息，包括 id 、 username 、 password 、 email 等字段。
type User struct {
	gorm.Model
	UserName string
	Password string
	Email    string `json:"Email,omitempty"`
	Phone    string `json:"Phone,omitempty"`
	Posts    []Post `json:"Posts,omitempty"`
}

func (u *User) SetPassword(pw string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}

func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

// posts 表：存储博客文章信息，包括 id 、 title 、 content 、 user_id （关联 users 表的 id）、
// created_at 、 updated_at 等字段。
type Post struct {
	gorm.Model
	Title    string
	Content  string
	UserID   uint
	User     *User     `gorm:"foreignKey:UserID" json:"User,omitempty"`
	Comments []Comment `json:"Comments,omitempty"`
}

// comments 表：存储文章评论信息，包括 id 、 content 、 user_id （关联 users 表的 id）、
// post_id （关联 posts 表的 id ）、 created_at 等字段。
type Comment struct {
	gorm.Model
	Content string
	UserID  uint
	PostID  uint
}
