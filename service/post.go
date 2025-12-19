package service

import (
	"github.com/xiaohan1995/Gin-blog/models"
	"gorm.io/gorm"
)

func CreatePost(post models.Post) *models.APIError {
	DB := models.DB
	err := DB.Create(&post).Error
	if err != nil {
		models.Log.Error("创建文章失败", err)
		return models.ErrInternalServer
	}
	models.Log.Info("新文章被创建:", post.ID)
	return nil
}

func GetPosts() ([]models.Post, int, *models.APIError) {
	var posts []models.Post
	DB := models.DB
	err := DB.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Omit("Password")
	}).Find(&posts).Error
	if err != nil {
		models.Log.Error("获取文章列表失败:", err)
		return posts, 0, models.ErrInternalServer
	}
	models.Log.Info("获取文章列表成功:", len(posts))
	return posts, len(posts), nil
}

func GetPost(id uint, userID interface{}) (models.Post, *models.APIError) {
	var post models.Post
	DB := models.DB
	err := DB.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Omit("Password")
	}).Where("id = ?", id).First(&post).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			models.Log.Warning("未找到文章ID:", id)
			return post, models.ErrPostNotFound
		}
		models.Log.Error("获取文章失败:", err)
		return post, models.ErrInternalServer
	}
	models.Log.Info("获取文章成功:", post.ID)
	return post, nil
}

func UpdatePost(id uint, userID interface{}, post models.Post) *models.APIError {
	DB := models.DB
	// 先检查文章是否存在
	var existingPost models.Post
	err := DB.Where("id = ?", id).First(&existingPost).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			models.Log.Warning("文章不存在:", id)
			return models.ErrPostNotFound
		}

		models.Log.Error("查询文章失败:", err)
		return models.ErrInternalServer
	}

	// 判断文章user_id是否等于当前用户
	if existingPost.UserID != userID {
		models.Log.Error("无权修改该文章:", existingPost.UserID, "Current user:", userID)
		return models.ErrForbidden
	}
	err = DB.Where("id = ? and user_id = ?", id, userID).Updates(&post).Error
	if err != nil {
		models.Log.Error("更新文章失败:", err)
		return models.ErrInternalServer
	}
	models.Log.Info("文章被更新:", post.ID)
	return nil
}

func DeletePost(id uint, userID interface{}) *models.APIError {
	DB := models.DB
	post := models.Post{}
	// 先检查文章是否存在
	var existingPost models.Post
	err := DB.Where("id = ?", id).First(&existingPost).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			models.Log.Warning("文章不存在:", id)
			return models.ErrPostNotFound
		}

		models.Log.Error("查询文章失败:", err)
		return models.ErrInternalServer
	}

	// 判断文章user_id是否等于当前用户
	if existingPost.UserID != userID {
		models.Log.Error("无权删除该文章:", existingPost.UserID, "Current user:", userID)
		return models.ErrForbidden
	}
	err = DB.Where("id = ? and user_id = ?", id, userID).Delete(&post).Error
	if err != nil {
		models.Log.Error("删除文章失败:", err)
	}
	models.Log.Info("删除文章成功:", id)
	return nil
}

// 获取文章的评论列表
func GetPostComments(id int) ([]models.Comment, *models.APIError) {
	var comments []models.Comment
	DB := models.DB
	err := DB.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,user_name")
	}).Where("post_id = ?", id).Find(&comments).Error
	if err != nil {
		models.Log.Error("获取文章评论失败:", err)
		return nil, models.ErrInternalServer
	}
	return comments, nil
}
