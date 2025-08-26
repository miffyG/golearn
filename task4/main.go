package main

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/gin-gonic/gin"
	"github.com/miffyG/golearn/task4/config"
	"github.com/miffyG/golearn/task4/db"
	"github.com/miffyG/golearn/task4/handler"
	"github.com/miffyG/golearn/task4/logger"
	"github.com/miffyG/golearn/task4/middleware"
	"github.com/miffyG/golearn/task4/models"
	"github.com/miffyG/golearn/task4/repository"
	"github.com/miffyG/golearn/task4/service"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	config.LoadConfig()

	if err := logger.Init(); err != nil {
		log.Fatalf("日志初始化失败: %v", err)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			_ = err
		}
	}()

	db.InitGormDb(config.GetDbConfig())

	if err := db.GormDb.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{}); err != nil {
		logger.Sugar.Fatalf("数据库自动迁移失败: %v", err)
	}
	// insertBlogTestData()

	userRepo := repository.NewUserRepository(db.GormDb)
	postRepo := repository.NewPostRepository(db.GormDb)
	commentRepo := repository.NewCommentRepository(db.GormDb)

	userService := service.NewUserService(userRepo)
	postService := service.NewPostService(postRepo)
	commentService := service.NewCommentService(commentRepo)

	userHandler := handler.NewAuthHandler(userService)
	postHandler := handler.NewPostHandler(postService)
	commentHandler := handler.NewCommentHandler(commentService)

	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	setupRoutes(r, userHandler, postHandler, commentHandler)
	logger.Sugar.Info("服务器启动，监听端口 8080")
	if err := r.Run(":8080"); err != nil {
		logger.Sugar.Fatalf("服务器启动失败: %v", err)
	}

}

func setupRoutes(r *gin.Engine, userHandler *handler.AuthHandler, postHandler *handler.PostHandler, commentHandler *handler.CommentHandler) {
	v1 := r.Group("/api/v1")
	{
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/register", userHandler.Register)
			authGroup.POST("/login", userHandler.Login)
		}

		v1.GET("/posts", postHandler.GetPosts)
		v1.GET("/posts/:post_id", postHandler.GetPostsById)

		protected := v1.Group("/")
		protected.Use(middleware.JwtAuthMiddleware())
		{
			protected.POST("/posts", postHandler.CreatePost)
			protected.PUT("/posts/:post_id", postHandler.UpdatePost)
			protected.DELETE("/posts/:post_id", postHandler.DeletePost)
			protected.POST("/posts/:post_id/comments", commentHandler.CreateComment)
		}

		v1.GET("/posts/:post_id/comments", commentHandler.GetCommentsByPost)
	}
}

// 如果没有数据则插入一些数据
func insertBlogTestData() {
	db := db.GormDb
	var userCount int64
	if err := db.Model(&models.User{}).Count(&userCount).Error; err != nil {
		fmt.Println("count error:", err)
		return
	}
	if userCount > 0 {
		fmt.Println("已有博客测试数据，无须插入")
		return
	}

	// 可调整 batch 大小
	const postBatch = 100
	const commentBatch = 500

	for i := 1; i <= 10; i++ {
		// 开事务，保证一组用户相关写要么都成功要么回滚
		tx := db.Begin()
		if tx.Error != nil {
			fmt.Println("begin tx err:", tx.Error)
			return
		}

		user := models.User{UserName: fmt.Sprintf("User%d", i), Email: fmt.Sprintf("user%d@example.com", i)}
		if err := tx.Create(&user).Error; err != nil {
			tx.Rollback()
			fmt.Println("插入用户数据失败:", err)
			return
		}

		// 先生成 posts 切片（不逐条写）
		numPosts := rand.Intn(15) + 1
		posts := make([]models.Post, 0, numPosts)
		for j := 1; j <= numPosts; j++ {
			posts = append(posts, models.Post{
				Title:   fmt.Sprintf("Post %d by User%d", j, i),
				Content: "This is a sample post.",
				UserID:  user.ID,
			})
		}
		// 批量写 posts（GORM 会回填 posts 中的 ID）
		if err := tx.CreateInBatches(&posts, postBatch).Error; err != nil {
			tx.Rollback()
			fmt.Println("插入文章数据失败:", err)
			return
		}

		// 生成 comments 切片（使用 posts 回填的 ID）
		comments := make([]models.Comment, 0, numPosts*5) // 估算容量
		for idx, p := range posts {
			numComments := rand.Intn(13) + 1
			for k := 1; k <= numComments; k++ {
				comments = append(comments, models.Comment{
					Content: fmt.Sprintf("Comment %d on Post %d by User%d", k, idx+1, i),
					PostID:  p.ID,
					UserID:  user.ID,
				})
			}
		}
		if len(comments) > 0 {
			if err := tx.CreateInBatches(&comments, commentBatch).Error; err != nil {
				tx.Rollback()
				fmt.Println("插入评论数据失败:", err)
				return
			}
		}

		// 更新用户文章数（可以直接 DB 更新避免再次读写）
		// if err := tx.Model(&user).Update("post_count", len(posts)).Error; err != nil {
		// 	tx.Rollback()
		// 	fmt.Println("更新用户文章数失败:", err)
		// 	return
		// }

		if err := tx.Commit().Error; err != nil {
			fmt.Println("commit err:", err)
			return
		}
	}

	fmt.Println("插入博客测试数据成功")
}
