package main

import (
	"fmt"
	"math/rand"

	"github.com/miffyG/golearn/task3/db"
	"gorm.io/gorm"
)

/** 要开发一个博客系统，有以下几个实体： BlogUser （用户）、 Post （文章）、 Comment （评论）
 *  使用Gorm定义 BlogUser 、 Post 和 Comment 模型，其中 BlogUser 与 Post 是一对多关系（一个用户可以发布多篇文章），
 *  Post 与 Comment 也是一对多关系（一篇文章可以有多个评论）
 */
type BlogUser struct {
	gorm.Model
	Name      string
	Email     string
	Posts     []Post
	PostCount int `gorm:"default:0"` // 文章数量统计字段
}

type Post struct {
	gorm.Model
	Title      string
	Content    string
	BlogUserID uint
	Status     string `gorm:"default:''"` // 文章状态，默认为空字符
	Comments   []Comment
}

type Comment struct {
	gorm.Model
	Content    string
	PostID     uint
	BlogUserID uint
}

// 创建这些模型对应的数据库表
func createBlogTables() {
	db := db.GormDb
	if err := db.AutoMigrate(&BlogUser{}, &Post{}, &Comment{}); err != nil {
		panic("failed to migrate blog tables: " + err.Error())
	}
	fmt.Println("创建博客相关User表、Post表、Comment表成功")
}

// 如果没有数据则插入一些数据
func insertBlogTestData() {
	/* db := db.GormDb
	var userCount int64
	db.Model(&BlogUser{}).Count(&userCount)
	if userCount > 0 {
		fmt.Println("已有博客测试数据，无须插入")
		return
	}
	for i := 1; i <= 10; i++ {
		user := BlogUser{Name: fmt.Sprintf("User %d", i), Email: fmt.Sprintf("user%d@example.com", i)}
		res := db.Create(&user)
		if res.Error != nil {
			fmt.Println("插入用户数据失败:", res.Error)
			return
		}
		for j := 1; j <= rand.Intn(15)+1; j++ {
			post := Post{Title: fmt.Sprintf("Post %d by User %d", j, i), Content: "This is a sample post.", BlogUserID: user.ID}
			res := db.Create(&post)
			if res.Error != nil {
				fmt.Println("插入文章数据失败:", res.Error)
				return
			}
			for k := 1; k <= rand.Intn(13)+1; k++ {
				comment := Comment{Content: fmt.Sprintf("Comment %d on Post %d by User %d", k, j, i), PostID: post.ID, BlogUserID: user.ID}
				res := db.Create(&comment)
				if res.Error != nil {
					fmt.Println("插入评论数据失败:", res.Error)
					return
				}
			}
		}
	}
	fmt.Println("插入博客测试数据成功") */
	db := db.GormDb
	var userCount int64
	if err := db.Model(&BlogUser{}).Count(&userCount).Error; err != nil {
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

		user := BlogUser{Name: fmt.Sprintf("User %d", i), Email: fmt.Sprintf("user%d@example.com", i)}
		if err := tx.Create(&user).Error; err != nil {
			tx.Rollback()
			fmt.Println("插入用户数据失败:", err)
			return
		}

		// 先生成 posts 切片（不逐条写）
		numPosts := rand.Intn(15) + 1
		posts := make([]Post, 0, numPosts)
		for j := 1; j <= numPosts; j++ {
			posts = append(posts, Post{
				Title:      fmt.Sprintf("Post %d by User %d", j, i),
				Content:    "This is a sample post.",
				BlogUserID: user.ID,
			})
		}
		// 批量写 posts（GORM 会回填 posts 中的 ID）
		if err := tx.CreateInBatches(&posts, postBatch).Error; err != nil {
			tx.Rollback()
			fmt.Println("插入文章数据失败:", err)
			return
		}

		// 生成 comments 切片（使用 posts 回填的 ID）
		comments := make([]Comment, 0, numPosts*5) // 估算容量
		for idx, p := range posts {
			numComments := rand.Intn(13) + 1
			for k := 1; k <= numComments; k++ {
				comments = append(comments, Comment{
					Content:    fmt.Sprintf("Comment %d on Post %d by User %d", k, idx+1, i),
					PostID:     p.ID,
					BlogUserID: user.ID,
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

// 查询某个用户发布的所有文章及其对应的评论信息
func getUserPostsAndComments(blogUserID uint) {
	var posts []Post
	db := db.GormDb
	if err := db.Preload("Comments").Where("blog_user_id = ?", blogUserID).Find(&posts).Error; err != nil {
		fmt.Println("查询用户文章失败:", err)
	} else {
		fmt.Println("查询到的用户文章及评论信息:")
		for _, post := range posts {
			fmt.Printf("文章ID: %d, 标题: %s, 内容: %s, 评论数量: %d\n",
				post.ID, post.Title, post.Content, len(post.Comments))
		}
	}
}

// 查询评论数量最多的文章信息
func getMostCommentedPost() {
	// 1) 取出最大评论数
	var maxCount int64
	err := db.GormDb.Model(&Comment{}).
		Select("COUNT(*) AS cnt").
		Group("post_id").
		Order("cnt DESC").
		Limit(1).
		Pluck("cnt", &maxCount).Error
	if err != nil {
		fmt.Println("获取最大评论数失败:", err)
		return
	}
	if maxCount == 0 {
		fmt.Println("没有评论记录")
		return
	}

	// 2) 查出所有评论数等于 maxCount 的 post_id
	var postIDs []uint
	err = db.GormDb.Model(&Comment{}).
		Select("post_id").
		Group("post_id").
		Having("COUNT(*) = ?", maxCount).
		Pluck("post_id", &postIDs).Error
	if err != nil {
		fmt.Println("查找拥有最大评论数的 post_id 失败:", err)
		return
	}
	if len(postIDs) == 0 {
		fmt.Println("没有找到任何文章")
		return
	}

	// 3) 一次性加载这些文章及其评论
	var posts []Post
	if err := db.GormDb.Preload("Comments").Find(&posts, postIDs).Error; err != nil {
		fmt.Println("根据 id 加载文章失败:", err)
		return
	}

	fmt.Printf("评论数最多的文章（评论数=%d），共 %d 篇：\n", maxCount, len(posts))
	for _, p := range posts {
		fmt.Printf("Post ID=%d, Title=%q, Comments=%d\n", p.ID, p.Title, len(p.Comments))
	}
}

// 为 Post 模型添加一个钩子函数，在文章创建时自动更新用户的文章数量统计字段
func (p *Post) AfterCreate(tx *gorm.DB) error {
	var user BlogUser
	fmt.Println("文章数量统计字段更新前是：", user.PostCount)
	if err := tx.Model(&BlogUser{}).Where("id = ?", p.BlogUserID).First(&user).Error; err != nil {
		fmt.Println("文章创建时更新用户的文章数量统计字段失败:", err)
		return err
	}
	user.PostCount++
	if err := tx.Save(&user).Error; err != nil {
		fmt.Println("文章创建时更新用户的文章数量统计字段失败:", err)
		return err
	} else {
		fmt.Println("文章创建时更新用户的文章数量统计字段成功:")
		fmt.Println("文章数量统计字段更新后是：", user.PostCount)
	}
	return nil
}

// 为 Comment 模型添加一个钩子函数，在评论删除时检查文章的评论数量，如果评论数量为 0，则更新文章的评论状态为 "无评论"
func (c *Comment) AfterDelete(tx *gorm.DB) error {
	var post Post
	if err := tx.Model(&Post{}).Where("id = ?", c.PostID).First(&post).Error; err != nil {
		fmt.Println("评论删除时查询文章失败:", err)
		return err
	}
	var comments []Comment
	if err := tx.Model(&Comment{}).Where("post_id = ?", c.PostID).Find(&comments).Error; err != nil {
		fmt.Println("评论删除时查询文章评论失败:", err)
		return err
	}
	if len(comments) == 0 {
		post.Status = "无评论"
		if err := tx.Save(&post).Error; err != nil {
			fmt.Println("评论删除时更新文章状态失败:", err)
			return err
		} else {
			fmt.Println("评论删除时更新文章状态成功:", post)
		}
	}
	return nil
}
