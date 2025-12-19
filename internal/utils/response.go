package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// PageData 分页数据结构
type PageData struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: 200,
		Msg:  "success",
		Data: data,
	})
}

// SuccessWithMsg 成功响应（带消息）
func SuccessWithMsg(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: 200,
		Msg:  msg,
		Data: data,
	})
}

// Error 错误响应
func Error(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: 500,
		Msg:  msg,
	})
}

// ErrorWithCode 错误响应（带状态码）
func ErrorWithCode(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
	})
}

// BadRequest 错误请求
func BadRequest(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: 400,
		Msg:  msg,
	})
}

// Unauthorized 未授权
func Unauthorized(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: 401,
		Msg:  msg,
	})
}

// Forbidden 禁止访问
func Forbidden(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: 403,
		Msg:  msg,
	})
}

// NotFound 未找到
func NotFound(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: 404,
		Msg:  msg,
	})
}

// PageSuccess 分页成功响应
func PageSuccess(c *gin.Context, list interface{}, total int64, page, size int) {
	c.JSON(http.StatusOK, Response{
		Code: 200,
		Msg:  "success",
		Data: PageData{
			List:  list,
			Total: total,
			Page:  page,
			Size:  size,
		},
	})
}
