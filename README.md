# 个人博客系统后台系统
使用 Go 语言结合 Gin 框架和 GORM 库开发一个个人博客系统的后端，实现博客文章的基本管理功能，包括文章的创建、读取、更新和删除（CRUD）操作，同时支持用户认证和简单的评论功能。

## 1.运行环境
  - Go 1.23.0
  - Gin 1.11.0
  - GORM 1.31.1

## 2.项目初始化
  1. 创建项目目录
  2. 初始化项目

    
    go mod init github.com/xiaohan1995/Gin-blog
    

  3. 创建main.go文件
  4. 安装GORM和GIN

    
    go get -u gorm.io/gorm

    go get -u gorm.io/driver/mysql

    go get -u github.com/gin-gonic/gin

    go mod tidy

    

## 3.项目结构
  ```
   Gin-blog
    |-models
    |-routers
    |-service
    |-statics
        |-css
        |-js
    |-templates
    |-main.go
    |-go.mod
    |-go.sum
    |-README.md

  ```
## 4.main.go
  ```go
    package main
    import (
        "github.com/xiaohan1995/Gin-blog/models"
        "github.com/xiaohan1995/Gin-blog/routers"

        "github.com/gin-gonic/gin"
    )
    func main() {
        // 初始化数据库
        models.InitDB()
        //初始化gin
        r := gin.Default()
        //设置静态资源和模版路径
        r.Static("/statics", "./statics")
        r.LoadHTMLGlob("templates/*")

        //初始化路由
        routers.InitRouter(r)
        r.Run(":8081")`
    }
  ```

## 5.数据库连接配置
  - 采用mysql数据库

  ```go
    //数据库连接初始化模型
    func InitDB() {
        dsn := "root:yourpassword@tcp(127.0.0.1:3306)/blog?charset=utf8mb4&parseTime=True&loc=Local"
        var err error
        DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
        if err != nil {
            log.Fatal("Failed to connect to database:", err)
        }
        DB.AutoMigrate(&User{}, &Post{}, &Comment{})
    }
```

## 6.启动服务
```
    go run main.go
```

## 7.测试
  - 创建用户：


