package routers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xiaohan1995/Gin-blog/middleware"
	"github.com/xiaohan1995/Gin-blog/models"
	"github.com/xiaohan1995/Gin-blog/service"
)

func InitRouter(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to Gin Blog API",
		})
	})

	//注册页面
	r.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", nil)
	})

	//注册api
	r.POST("/api/register", func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBind(&user); err != nil {
			models.Log.Warning("注册请求无效:", err)
			c.JSON(models.ErrInvalidRequest.Code, models.ErrInvalidRequest)
			return
		}

		if apiErr := service.RegisterUser(user); apiErr != nil {
			c.JSON(apiErr.Code, apiErr)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "注册成功",
		})
	})

	//登录页面
	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	//登录api
	r.POST("/api/login", func(c *gin.Context) {
		// 获取用户名和密码
		username := c.PostForm("username")
		password := c.PostForm("password")

		// 如果没有通过form获取到，则尝试通过JSON获取
		if username == "" && password == "" {
			var loginReq struct {
				Username string `json:"username"`
				Password string `json:"password"`
			}
			if err := c.ShouldBindJSON(&loginReq); err == nil {
				username = loginReq.Username
				password = loginReq.Password
			}
		}
		res, apiErr := service.LoginUser(username, password)
		if apiErr != nil {
			c.JSON(apiErr.Code, apiErr)
			return
		}
		c.JSON(http.StatusOK, res)

	})

	// //管理主页
	r.GET("/admin", func(c *gin.Context) {
		c.HTML(http.StatusOK, "admin.html", nil)
	})

	//用户列表
	r.GET("/users", func(c *gin.Context) {
		c.HTML(http.StatusOK, "users.html", nil)
	})

	//文章列表
	r.GET("/posts", func(c *gin.Context) {
		c.HTML(http.StatusOK, "posts.html", nil)
	})

	//评论列表
	r.GET("/comments", func(c *gin.Context) {
		c.HTML(http.StatusOK, "comments.html", nil)
	})

	//文章详情页面
	r.GET("/post-detail/:id", func(c *gin.Context) {
		postID := c.Param("id")
		if postID == "" {
			models.Log.Error("文章ID为空")
			c.JSON(models.ErrInvalidRequest.Code, models.ErrInvalidRequest)
		}
		c.HTML(http.StatusOK, "post-detail.html", gin.H{"postID": postID})

	})

	// 受保护的路由示例
	protected := r.Group("/api/protected")
	protected.Use(middleware.AuthMiddleware)
	{
		//用户列表
		protected.GET("/users", func(c *gin.Context) {
			users, err := service.GetUsers()
			if err != nil {
				c.JSON(models.ErrDatabaseConnection.Code, models.ErrDatabaseConnection)
				return
			}
			c.JSON(http.StatusOK, users)
		})
		//获取文章列表
		protected.GET("/posts", func(c *gin.Context) {
			posts, count, apiErr := service.GetPosts()
			if apiErr != nil {
				c.JSON(http.StatusOK, gin.H{
					"message": apiErr.Message,
					"code":    apiErr.Code,
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"message": "获取文章列表成功",
				"data":    posts,
				"count":   count,
			})
		})
		//添加文章
		protected.POST("/posts", func(c *gin.Context) {
			var postReq struct {
				Title   string `json:"title" binding:"required"`
				Content string `json:"content" binding:"required"`
			}
			if err := c.ShouldBind(&postReq); err != nil {
				models.Log.Error(err.Error())
				c.JSON(models.ErrInvalidRequest.Code, models.ErrInvalidRequest)
				return
			}
			// 从上下文中获取用户ID
			UserID, exists := c.Get("user_id")
			if !exists {
				models.Log.Error("无法获取用户信息")
				c.JSON(models.ErrInvalidRequest.Code, models.ErrInvalidRequest)
				return
			}
			post := models.Post{
				Title:   postReq.Title,
				Content: postReq.Content,
				UserID:  UserID.(uint),
			}
			if apiErr := service.CreatePost(post); apiErr != nil {
				c.JSON(apiErr.Code, apiErr)
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"message": "文章创建成功",
			})
		})
		//获取文章信息
		protected.GET("/post/:id", func(c *gin.Context) {
			postID := c.Param("id")
			if postID == "" {
				models.Log.Error("文章ID为空")
				c.JSON(models.ErrInvalidRequest.Code, models.ErrInvalidRequest)
			}
			postIDInt, err := strconv.Atoi(postID)
			if err != nil {
				models.Log.Error("文章ID转换失败:", err)
				c.JSON(models.ErrInvalidRequest.Code, models.ErrInvalidRequest)
			}
			// 从上下文中获取用户ID
			UserID, exists := c.Get("user_id")
			if !exists {
				models.Log.Error("无法获取用户信息")
				c.JSON(models.ErrUnauthorized.Code, models.ErrUnauthorized)
				return
			}
			post, apiErr := service.GetPost(uint(postIDInt), UserID)
			if apiErr != nil {
				c.JSON(apiErr.Code, apiErr)
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"message": "获取文章成功",
				"data":    post,
			})
		})

		//更新文章
		protected.PUT("/post/:id", func(c *gin.Context) {
			var postReq struct {
				Title   string `json:"title" binding:"required"`
				Content string `json:"content" binding:"required"`
			}
			if err := c.ShouldBind(&postReq); err != nil {
				models.Log.Error(err)
				c.JSON(models.ErrInvalidRequest.Code, models.ErrInvalidRequest)
				return
			}
			postID := c.Param("id")
			postIDInt, err := strconv.Atoi(postID)
			if err != nil {
				models.Log.Error("文章ID错误:", err)
				c.JSON(models.ErrInvalidRequest.Code, models.ErrInvalidRequest)
				return
			}
			UserID, exists := c.Get("user_id")
			if !exists {
				models.Log.Error("无法获取用户信息")
				c.JSON(models.ErrUnauthorized.Code, models.ErrUnauthorized)
				return
			}
			post := models.Post{
				Title:   postReq.Title,
				Content: postReq.Content,
				UserID:  UserID.(uint),
			}
			apiErr := service.UpdatePost(uint(postIDInt), UserID, post)
			if apiErr != nil {
				c.JSON(apiErr.Code, apiErr)
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"message": "更新成功",
			})

		})
		//删除文章
		protected.DELETE("/post/:id", func(c *gin.Context) {
			postID := c.Param("id")
			if postID == "" {
				models.Log.Error("文章ID为空")
				c.JSON(models.ErrInvalidRequest.Code, models.ErrInvalidRequest)
			}
			postIDInt, err := strconv.Atoi(postID)
			if err != nil {
				models.Log.Error("文章ID错误:", err)
				c.JSON(models.ErrInvalidRequest.Code, models.ErrInvalidRequest)
			}
			// 从上下文中获取用户ID
			UserID, exists := c.Get("user_id")
			if !exists {
				models.Log.Error("无法获取用户信息")
				c.JSON(models.ErrUnauthorized.Code, models.ErrUnauthorized)
				return
			}
			apiErr := service.DeletePost(uint(postIDInt), UserID)
			if apiErr != nil {
				c.JSON(apiErr.Code, apiErr)
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"message": "删除成功",
			})
		})

		//获取文章的评论列表
		protected.GET("/post/:id/comments", func(c *gin.Context) {
			postID := c.Param("id")
			if postID == "" {
				models.Log.Error("文章ID为空")
				c.JSON(models.ErrInvalidRequest.Code, models.ErrInvalidRequest)
			}
			postIDInt, err := strconv.Atoi(postID)
			if err != nil {
				models.Log.Error("文章ID错误:", err)
				c.JSON(models.ErrInvalidRequest.Code, models.ErrInvalidRequest)
			}
			res, apiErr := service.GetPostComments(postIDInt)
			if apiErr != nil {
				c.JSON(apiErr.Code, apiErr)
			}
			c.JSON(http.StatusOK, gin.H{
				"message": "获取评论列表成功",
				"data":    res,
			})

		})

		//创建评论
		protected.POST("/comments", func(c *gin.Context) {
			var commentReq struct {
				Content string `json:"content" binding:"required"`
				PostID  uint   `json:"post_id" binding:"required"`
			}

			err := c.ShouldBindJSON(&commentReq)
			if err != nil {
				models.Log.Error("参数错误:", err)
				c.JSON(models.ErrInvalidRequest.Code, models.ErrInvalidRequest)
				return
			}
			// 从上下文中获取用户ID
			UserID, exists := c.Get("user_id")
			if !exists {
				models.Log.Error("无法获取用户信息")
				c.JSON(models.ErrInvalidRequest.Code, models.ErrInvalidRequest)
				return
			}
			comment := models.Comment{
				Content: commentReq.Content,
				PostID:  commentReq.PostID,
				UserID:  UserID.(uint),
			}
			apiErr := service.CreateComment(comment)
			if apiErr != nil {
				models.Log.Error("创建评论失败:", apiErr)
				c.JSON(apiErr.Code, apiErr)
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "评论成功"})
		})

		//获取评论列表
		protected.GET("/comments", func(c *gin.Context) {
			comments, apiErr := service.GetComments()
			if apiErr != nil {
				models.Log.Error("获取评论列表失败:", apiErr)
				c.JSON(apiErr.Code, apiErr)
				return
			}
			c.JSON(http.StatusOK, gin.H{"data": comments, "message": "获取评论列表成功"})
		})

		//用户信息
		protected.GET("/profile", func(c *gin.Context) {
			userID, _ := c.Get("user_id")
			username, _ := c.Get("username")
			c.JSON(http.StatusOK, gin.H{
				"user_id":  userID,
				"username": username,
				"message":  "这是受保护的资源",
			})
		})
	}
}
