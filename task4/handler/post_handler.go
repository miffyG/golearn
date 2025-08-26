package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/miffyG/golearn/task4/logger"
	"github.com/miffyG/golearn/task4/models"
	"github.com/miffyG/golearn/task4/service"
	"gorm.io/gorm"
)

type PostHandler struct {
	service *service.PostService
}

func NewPostHandler(s *service.PostService) *PostHandler {
	return &PostHandler{service: s}
}

type CreatePostRequest struct {
	Title   string `json:"title" form:"title" binding:"required"`
	Content string `json:"content" form:"content" binding:"required"`
}

// @Summary 创建帖子
// @Description 创建帖子接口
// @Tags posts
// @Accept json
// @Produce json
// @Param post body CreatePostRequest true "帖子信息"
// @Success 200 {object} models.Response "{"code":200,"data":{"post_id":1},"msg":"创建帖子成功"}"
// @Failure 400 {object} models.ErrorResponse "{"code":400,"msg":"参数错误"}"
// @Failure 500 {object} models.ErrorResponse "{"code":500,"msg":"创建帖子失败"}"
// @Router /posts [post]
func (h *PostHandler) CreatePost(c *gin.Context) {
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Errorf("参数错误: %v", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    400,
			Message: "参数错误",
		})
		return
	}

	userId := c.GetUint("user_id")
	post := models.Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userId,
	}
	if err := h.service.Create(&post); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    500,
			Message: "创建帖子失败",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "创建帖子成功",
		Data:    post,
	})
}

// @Summary 获取帖子列表
// @Description 获取帖子列表接口
// @Tags posts
// @Accept json
// @Produce json
// @Success 200 {object} models.Response "{"code":200,"data":[{"id":1,"title":"标题","user_id":1}],"msg":"获取帖子成功"}"
// @Failure 500 {object} models.ErrorResponse "{"code":500,"msg":"获取帖子失败"}"
// @Router /posts [get]
func (h *PostHandler) GetPosts(c *gin.Context) {
	posts, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    500,
			Message: "获取帖子失败",
		})
		return
	}

	postResponses := make([]models.PostResponse, len(posts))
	for i, post := range posts {
		postResponses[i] = models.PostResponse{
			Title:   post.Title,
			Content: post.Content,
			UserID:  post.UserID,
			User: &models.UserResponse{
				ID:       post.User.ID,
				Username: post.User.UserName,
				Email:    post.User.Email,
				Phone:    post.User.Phone,
			},
		}
		for _, comment := range post.Comments {
			postResponses[i].Comments = append(postResponses[i].Comments, models.CommentResponse{
				Content: comment.Content,
				UserID:  comment.UserID,
				PostID:  comment.PostID,
			})
		}
	}
	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "获取帖子成功",
		Data:    postResponses,
	})
}

// @Summary 获取帖子详情
// @Description 获取帖子详情接口
// @Tags posts
// @Accept json
// @Produce json
// @Param post_id path int true "帖子ID"
// @Success 200 {object} models.Response "{"code":200,"data":{"id":1,"title":"标题","content":"内容","user_id":1},"msg":"获取帖子成功"}"
// @Failure 400 {object} models.ErrorResponse "{"code":400,"msg":"参数错误"}"
// @Failure 404 {object} models.ErrorResponse "{"code":404,"msg":"帖子未找到"}"
// @Failure 500 {object} models.ErrorResponse "{"code":500,"msg":"获取帖子失败"}"
// @Router /posts/{post_id} [get]
func (h *PostHandler) GetPostsById(c *gin.Context) {
	postIdStr := c.Param("post_id")
	var postId uint
	_, errConv := fmt.Sscanf(postIdStr, "%d", &postId)
	if errConv != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    400,
			Message: "参数错误",
		})
		return
	}
	p, err := h.service.GetByID(postId)
	if err != nil {
		if err.Error() != "record not found" {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Code:    500,
				Message: "获取帖子失败",
			})
			return
		}
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Code:    404,
			Message: "帖子未找到",
		})
		return
	}

	postRes := models.PostResponse{
		Title:   p.Title,
		Content: p.Content,
		UserID:  p.UserID,
		User: &models.UserResponse{
			ID:       p.User.ID,
			Username: p.User.UserName,
			Email:    p.User.Email,
			Phone:    p.User.Phone,
		},
	}
	for _, comment := range p.Comments {
		postRes.Comments = append(postRes.Comments, models.CommentResponse{
			Content: comment.Content,
		})
	}
	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "获取帖子成功",
		Data:    postRes,
	})
}

// @Summary 更新帖子
// @Description 更新帖子接口
// @Tags posts
// @Accept json
// @Produce json
// @Param post_id path int true "帖子ID"
// @Param post body CreatePostRequest true "帖子信息"
// @Success 200 {object} models.Response "{"code":200,"data":{"id":1,"title":"标题","content":"内容","user_id":1},"msg":"更新帖子成功"}"
// @Failure 400 {object} models.ErrorResponse "{"code":400,"msg":"参数错误"}"
// @Failure 403 {object} models.ErrorResponse "{"code":403,"msg":"没有权限"}"
// @Failure 404 {object} models.ErrorResponse "{"code":404,"msg":"帖子未找到"}"
// @Failure 500 {object} models.ErrorResponse "{"code":500,"msg":"更新帖子失败"}"
// @Router /posts/{post_id} [put]
func (h *PostHandler) UpdatePost(c *gin.Context) {
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    400,
			Message: "参数错误",
		})
		return
	}
	postIdStr := c.Param("post_id")
	var postId uint
	_, errConv := fmt.Sscanf(postIdStr, "%d", &postId)
	if errConv != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    400,
			Message: "参数错误",
		})
		return
	}
	userId := c.GetUint("user_id")
	post := &models.Post{
		Title:   req.Title,
		Content: req.Content,
	}
	post.ID = postId
	if err := h.service.Update(userId, post); err != nil {
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusForbidden, models.ErrorResponse{
				Code:    403,
				Message: "没有权限",
			})
			return
		} else if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Code:    404,
				Message: "帖子未找到",
			})
			return
		} else {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Code:    500,
				Message: "更新帖子失败",
			})
			return
		}
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "更新帖子成功",
		Data:    post,
	})
}

// @Summary 删除帖子
// @Description 删除帖子接口
// @Tags posts
// @Accept json
// @Produce json
// @Param post_id path int true "帖子ID"
// @Success 200 {object} models.Response "{"code":200,"msg":"删除帖子成功"}"
// @Failure 400 {object} models.ErrorResponse "{"code":400,"msg":"参数错误"}"
// @Failure 403 {object} models.ErrorResponse "{"code":403,"msg":"没有权限"}"
// @Failure 404 {object} models.ErrorResponse "{"code":404,"msg":"帖子未找到"}"
// @Failure 500 {object} models.ErrorResponse "{"code":500,"msg":"删除帖子失败"}"
// @Router /posts/{post_id} [delete]
func (h *PostHandler) DeletePost(c *gin.Context) {
	postIdStr := c.Param("post_id")
	var postId uint
	_, errConv := fmt.Sscanf(postIdStr, "%d", &postId)
	if errConv != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    400,
			Message: "参数错误",
		})
		return
	}
	userId := c.GetUint("user_id")
	if err := h.service.Delete(userId, postId); err != nil {
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusForbidden, models.ErrorResponse{
				Code:    403,
				Message: "没有权限",
			})
			return
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Code:    404,
				Message: "帖子未找到",
			})
			return
		} else {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Code:    500,
				Message: "删除帖子失败",
			})
			return
		}
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    200,
		Message: "删除帖子成功",
	})
}
