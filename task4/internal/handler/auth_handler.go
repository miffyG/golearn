package handler

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/miffyG/golearn/task4/internal/models/dto"
	"github.com/miffyG/golearn/task4/internal/models/entity"
	"github.com/miffyG/golearn/task4/internal/service"
	"github.com/miffyG/golearn/task4/pkg/logger"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	UserService *service.UserService
}

func NewAuthHandler(s *service.UserService) *AuthHandler {
	return &AuthHandler{
		UserService: s,
	}
}

type RegisterRequest struct {
	Username string `json:"username" form:"username" binding:"required,min=3,max=26,alphanum"`
	Password string `json:"password" form:"password" binding:"required,min=6,max=26"`
	Email    string `json:"email" form:"email" binding:"required,email"`
	Phone    string `json:"phone" form:"phone" binding:"phone"`
}

var phone validator.Func = func(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	if phone == "" {
		return true // 空值通过
	}
	return regexp.MustCompile(`^(\+86)?1[3-9]\d{9}$`).MatchString(phone)
}

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("phone", phone)
	}
}

// @Summary 用户注册
// @Description 用户注册接口
// @Tags auth
// @Accept json
// @Produce json
// @Param user body RegisterRequest true "用户信息"
// @Success 200 {object} dto.Response "{"code":200,"data":{"user_id":1},"msg":"注册成功"}"
// @Failure 400 {object} dto.ErrorResponse "{"code":400,"msg":"参数错误或密码设置失败"}"
// @Failure 500 {object} dto.ErrorResponse "{"code":500,"msg":"注册失败"}"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBind(&req); err != nil {
		logger.Sugar.Errorf("参数错误: %v", err)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    400,
			Message: "参数错误或密码设置失败",
		})
		return
	}

	user := entity.User{
		UserName: req.Username,
		Password: req.Password,
		Email:    req.Email,
	}
	if err := h.UserService.Register(&user); err != nil {
		logger.Sugar.Errorf("用户注册失败: %v", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    500,
			Message: "注册失败",
		})
		return
	}

	logger.Sugar.Infof("用户注册成功: id %s name: %s", user.UserName, user.ID)
	c.JSON(http.StatusOK, dto.Response{
		Code:    200,
		Message: "注册成功",
		Data:    map[string]interface{}{"user_id": user.ID},
	})
}

type LoginRequest struct {
	Username string `json:"username" form:"username" binding:"required,min=3,max=26,alphanum"`
	Password string `json:"password" form:"password" binding:"required,min=6,max=26"`
}

// @Summary 用户登录
// @Description 用户登录接口
// @Tags auth
// @Accept json
// @Produce json
// @Param user body LoginRequest true "用户信息"
// @Success 200 {object} dto.Response "{"code":200,"data":{"token":"xxx"},"msg":"登录成功"}"
// @Failure 400 {object} dto.ErrorResponse "{"code":400,"msg":"参数错误"}"
// @Failure 401 {object} dto.ErrorResponse "{"code":401,"msg":"无效的凭证"}"
// @Failure 500 {object} dto.ErrorResponse "{"code":500,"msg":"登录失败"}"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBind(&req); err != nil {
		logger.Sugar.Errorf("参数错误: %v", err)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    400,
			Message: "参数错误",
		})
		return
	}

	token, user, err := h.UserService.Login(req.Username, req.Password)

	if err != nil {
		logger.Sugar.Errorf("用户登录失败: %v", err)
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) || errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Code:    401,
				Message: "用户不存在或密码错误",
			})
		} else {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Code:    500,
				Message: "登录失败",
			})
		}
		return
	}
	if user == nil || token == "" {
		logger.Sugar.Warnf("用户或token为空！")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    500,
			Message: "登录失败",
		})
		return
	}

	logger.Sugar.Infof("用户登录成功: id %s name: %s", user.UserName, user.ID)
	c.JSON(http.StatusOK, dto.Response{
		Code:    200,
		Message: "登录成功",
		Data:    map[string]interface{}{"token": token},
	})
}
