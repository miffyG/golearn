package models

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
}

type PostResponse struct {
	Title    string
	Content  string
	UserID   uint
	User     *UserResponse     `json:"User,omitempty"`
	Comments []CommentResponse `json:"Comments,omitempty"`
}

type CommentResponse struct {
	Content string
	UserID  uint `json:"UserID,omitempty"`
	PostID  uint `json:"PostID,omitempty"`
}
