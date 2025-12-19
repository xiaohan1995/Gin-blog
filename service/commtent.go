package service

import (
	"github.com/xiaohan1995/Gin-blog/models"
	"gorm.io/gorm"
)

func CreateComment(comment models.Comment) *models.APIError {
	DB := models.DB
	if err := DB.Create(&comment).Error; err != nil {
		models.Log.Error("发布评论失败", err)
		return models.ErrPostCreated
	}
	return nil
}

func GetComments() ([]models.Comment, *models.APIError) {
	DB := models.DB
	var comments []models.Comment
	if err := DB.Preload("Post").Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Omit("password")
	}).Find(&comments).Error; err != nil {
		models.Log.Error("获取评论失败", err)
		return nil, models.ErrInternalServer
	}
	return comments, nil
}
