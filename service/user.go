package service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/xiaohan1995/Gin-blog/models"
	"golang.org/x/crypto/bcrypt"
)

type UserResponse struct {
	UserID    uint      `json:"id"`
	UserName  string    `json:"user_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// RegisterUser 用户注册
func RegisterUser(user models.User) *models.APIError {
	// 参数校验
	if user.UserName == "" || user.Email == "" || user.Password == "" {
		models.Log.Warning("用户名、邮箱、密码不能为空")
		return models.ErrInvalidRequest
	}
	// 用户名、邮箱是否注册过
	var existingUser models.User
	res := models.DB.Where("user_name = ? OR email =?", user.UserName, user.Email).First(&existingUser)
	if res.Error == nil {
		models.Log.Warning("用户名、邮箱已存在")
		return &models.APIError{
			Code:    409,
			Message: "用户名、邮箱已存在",
			Details: "用户名、邮箱已存在",
		}
	}
	//加密密码
	user.Password = EncryptPassword(user.Password)

	if err := models.DB.Create(&user).Error; err != nil {
		models.Log.Error("创建新用户失败", err)
		return models.ErrInternalServer
	}
	models.Log.Info("新用户被创建:", user.UserName)
	return nil
}

func LoginUser(username, password string) (map[string]interface{}, *models.APIError) {
	//查找用户
	var user models.User
	res := models.DB.Where("user_name = ?", username).First(&user)
	if res.Error != nil {
		models.Log.Error("用户不存在:", username)
		return nil, models.ErrUserNotFound
	}
	//验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		models.Log.Error("密码错误:", username)
		return nil, models.ErrInvalidCredentials
	}

	//生成JWT token
	token, err := GenerateJWT(user.ID, user.UserName)
	if err != nil {
		models.Log.Error("生成JWT失败:", err)
		return nil, models.ErrInternalServer
	}
	UserResponse := UserResponse{
		UserID:    user.ID,
		UserName:  user.UserName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
	return map[string]interface{}{
		"state":   0,
		"message": "登录成功",
		"data":    UserResponse,
		"token":   token,
	}, nil
}

// 获取用户列表
func GetUsers() ([]models.User, *models.APIError) {
	var users []models.User
	if err := models.DB.Omit("password").Find(&users).Error; err != nil {
		models.Log.Error("获取用户列表失败:", err)
		return nil, models.ErrInternalServer
	}
	return users, nil
}

// JWTClaims 自定义声明结构体
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateJWT 生成JWT Token
func GenerateJWT(userID uint, username string) (string, error) {
	// 设置过期时间为24小时
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &JWTClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "user_token",
		},
	}
	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥签名token，实际项目中应该从环境变量读取
	tokenString, err := token.SignedString([]byte("your_secret_key_here"))
	if err != nil {
		models.Log.Error("生成JWT Token失败：", err)
		return "", err
	}
	return tokenString, err
}

// ParseJWT 解析JWT Token
func ParseJWT(tokenString string) (*JWTClaims, error) {
	claims := &JWTClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("your_secret_key_here"), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		models.Log.Warning("无效的JWT Token")
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// EncryptPassword 加密密码
func EncryptPassword(p string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		models.Log.Error("密码加密失败")
		panic(err)
	}
	return string(hashedPassword)
}
