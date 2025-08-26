package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/miffyG/golearn/task4/config"
	"github.com/miffyG/golearn/task4/models"
	"github.com/miffyG/golearn/task4/utils"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "未授权",
			})
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "未授权",
			})
			return
		}

		tokenStr := parts[1]
		claims, err := utils.ParseJWTToken(config.GetSecretConfig().JwtSecret, tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "token无效",
			})
			return
		}
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}
