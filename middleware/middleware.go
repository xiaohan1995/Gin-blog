package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xiaohan1995/Gin-blog/service"
)

// 添加JWT中间件示例
func AuthMiddleware(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少访问令牌"})
		c.Abort()
		return
	}

	// 解析Bearer token格式
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	claims, err := service.ParseJWT(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
		c.Abort()
		return
	}

	// 将用户信息存储在上下文中供后续使用
	c.Set("user_id", claims.UserID)
	c.Set("username", claims.Username)
	c.Next()
}
