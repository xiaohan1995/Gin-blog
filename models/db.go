package models

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName string `gorm:"unique;not null" json:"username" binding:"required,min=3"`
	Email    string `gorm:"not null" json:"email" binding:"required,email"`
	Password string `gorm:"not null" json:"password" binding:"required,min=6"`
}

type Post struct {
	gorm.Model
	Title   string `gorm:"not null"`
	Content string `gorm:"not null"`
	UserID  uint
	User    User
}

type Comment struct {
	gorm.Model
	Content string `gorm:"not null"`
	UserID  uint
	PostID  uint
	User    User
	Post    Post
}

var DB *gorm.DB

func InitDB() {
	dsn := "root:password123@tcp(127.0.0.1:3306)/blog?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		Log.Error("数据库链接错误:", err)
	}
	Log.Info("数据库链接成功")
	DB.AutoMigrate(&User{}, &Post{}, &Comment{})
	Log.Info("数据库迁移成功")
}
