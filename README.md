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
   . ├── main.go # 项目入口文件
     ├── go.mod # Go 模块定义文件
     ├── go.sum # Go 模块校验文件
     ├── README.md # 项目说明文档
     ├── .gitignore # Git 忽略文件配置
     ├── config/ # 配置文件目录
     ├── logs/ # 日志文件目录
     │ ├── error.log # 错误日志
     │ ├── info.log # 信息日志
     │ └── warning.log # 警告日志
     ├── middleware/ # 中间件目录
     │ └── middleware.go # 自定义中间件
     ├── models/ # 数据模型目录
     │ ├── db.go # 数据库连接和模型定义
     │ ├── error.go # 错误处理定义
     │ └── log.go # 日志模块
     ├── routers/ # 路由配置目录
     │ └── router.go # 路由定义
     ├── service/ # 业务逻辑目录
     │ ├── commtent.go # 评论相关服务
     │ ├── post.go # 文章相关服务
     │ └── user.go # 用户相关服务
     ├── statics/ # 静态资源目录
     │ ├── css/ # CSS 样式文件
     │ │ └── admin.css # 后台样式
     │ └── js/ # JavaScript 脚本文件
     │ └── ajax.js # AJAX 相关脚本
     ├── templates/ # HTML 模板目录
     │ ├── admin.html # 后台管理页面
     │ ├── comments.html # 评论展示页面
     │ ├── login.html # 登录页面
     │ ├── post-detail.html # 文章详情页面
     │ ├── posts.html # 文章列表页面
     │ ├── register.html # 注册页面
     │ └── users.html # 用户列表页


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
  - 注册用户：


