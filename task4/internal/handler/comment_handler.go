package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/miffyG/golearn/task4/internal/models/dto"
	"github.com/miffyG/golearn/task4/internal/models/entity"
	"github.com/miffyG/golearn/task4/internal/service"
)

type CommentHandler struct {
	service *service.CommentService
}

func NewCommentHandler(s *service.CommentService) *CommentHandler {
	return &CommentHandler{service: s}
}

type CreateCommentRequest struct {
	Content string `json:"content" form:"content" binding:"required"`
}

// @Summary 创建评论
// @Description 创建评论接口
// @Tags comments
// @Accept json
// @Produce json
// @Param post_id path int true "帖子ID"
// @Param comment body CreateCommentRequest true "评论信息"
// @Success 200 {object} dto.Response "{"code":200,"data":{"id":1,"content":"内容","user_id":1,"post_id":1},"msg":"创建评论成功"}"
// @Failure 400 {object} dto.ErrorResponse "{"code":400,"msg":"参数错误"}"
// @Failure 500 {object} dto.ErrorResponse "{"code":500,"msg":"创建评论失败"}"
// @Router /posts/{post_id}/comments [post]
func (h *CommentHandler) CreateComment(c *gin.Context) {
	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    400,
			Message: "参数错误",
		})
		return
	}

	pidStr := c.Param("post_id")
	var postId uint
	if _, err := fmt.Sscanf(pidStr, "%d", &postId); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    400,
			Message: "参数错误",
		})
		return
	}
	userId := c.GetUint("user_id")
	comment := entity.Comment{
		Content: req.Content,
		UserID:  userId,
		PostID:  postId,
	}
	if err := h.service.Create(&comment); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    500,
			Message: "创建评论失败",
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    200,
		Message: "创建评论成功",
		Data:    comment,
	})
}

// @Summary 获取帖子评论
// @Description 获取帖子评论接口
// @Tags comments
// @Accept json
// @Produce json
// @Param post_id path int true "帖子ID"
// @Success 200 {object} dto.Response "{"code":200,"data":[{"id":1,"content":"内容","user_id":1,"post_id":1}],"msg":""}"
// @Failure 500 {object} dto.ErrorResponse "{"code":500,"msg":"获取评论失败"}"
// @Router /posts/{post_id}/comments [get]
func (h *CommentHandler) GetCommentsByPost(c *gin.Context) {
	pidStr := c.Param("post_id")
	var postId uint
	if _, err := fmt.Sscanf(pidStr, "%d", &postId); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "参数错误",
		})
		return
	}
	comments, err := h.service.GetByPostId(postId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "获取评论失败",
		})
		return
	}

	var commentRes []dto.CommentResponse
	for _, comment := range comments {
		commentRes = append(commentRes, dto.CommentResponse{
			Content: comment.Content,
			UserID:  comment.UserID,
			PostID:  comment.PostID,
		})
	}
	c.JSON(http.StatusOK, dto.Response{
		Code:    http.StatusOK,
		Message: "",
		Data:    commentRes,
	})
}
