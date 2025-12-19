package models

import (
	"net/http"
)

// APIError 定义API错误结构
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Error 实现error接口
func (e *APIError) Error() string {
	return e.Message
}

// 错误类型常量
var (
	ErrDatabaseConnection = &APIError{
		Code:    http.StatusInternalServerError, // 500
		Message: "数据库连接失败",
		Details: "无法连接到数据库，请稍后重试",
	}

	ErrUserNotFound = &APIError{
		Code:    http.StatusNotFound, //404
		Message: "用户不存在",
		Details: "指定的用户未找到",
	}

	ErrInvalidCredentials = &APIError{
		Code:    http.StatusUnauthorized, //401
		Message: "用户名或密码错误",
		Details: "用户名或密码错误",
	}

	ErrUnauthorized = &APIError{
		Code:    http.StatusUnauthorized, //401
		Message: "未授权访问",
		Details: "需要有效的身份验证令牌",
	}

	ErrForbidden = &APIError{
		Code:    http.StatusForbidden, //403
		Message: "您没有权限执行此操作",
		Details: "您没有权限执行此操作",
	}

	ErrPostNotFound = &APIError{
		Code:    http.StatusNotFound, //404
		Message: "资源不存在",
		Details: "指定的资源未找到",
	}

	ErrPostCreated = &APIError{
		Code:    http.StatusBadRequest, //400
		Message: "内容发布失败",
		Details: "请检查您的请求参数",
	}

	ErrInvalidRequest = &APIError{
		Code:    http.StatusBadRequest, //400
		Message: "请求参数无效",
		Details: "请检查您的请求参数",
	}

	ErrInternalServer = &APIError{
		Code:    http.StatusInternalServerError, //500
		Message: "服务器内部错误",
		Details: "服务器遇到意外情况，无法完成请求",
	}
)
