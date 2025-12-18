package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/xiaoxin/blog-backend/internal/services"
	"github.com/xiaoxin/blog-backend/internal/utils"
)

// UserController 用户控制器
type UserController struct {
	userService *services.UserService
}

// NewUserController 创建用户控制器实例
func NewUserController() *UserController {
	return &UserController{
		userService: services.NewUserService(),
	}
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	Email    string `json:"email" binding:"omitempty,email"`
	Nickname string `json:"nickname" binding:"max=50"`
}

// Register 用户注册
func (ctrl *UserController) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	user, err := ctrl.userService.Register(req.Username, req.Password, req.Email, req.Nickname)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.SuccessWithMsg(c, "注册成功", gin.H{
		"id":       user.ID,
		"username": user.Username,
	})
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login 用户登录
func (ctrl *UserController) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	token, refreshToken, err := ctrl.userService.Login(req.Username, req.Password)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"token":         token,
		"refresh_token": refreshToken,
	})
}

// GetProfile 获取当前用户信息
func (ctrl *UserController) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "未授权")
		return
	}

	user, err := ctrl.userService.GetUserByID(userID.(uint))
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Success(c, user)
}

// UpdateProfileRequest 更新用户信息请求
type UpdateProfileRequest struct {
	Nickname string `json:"nickname" binding:"max=50"`
	Email    string `json:"email" binding:"omitempty,email"`
	Avatar   string `json:"avatar" binding:"max=255"`
}

// UpdateProfile 更新用户信息
func (ctrl *UserController) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "未授权")
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := ctrl.userService.UpdateUser(userID.(uint), req.Nickname, req.Email, req.Avatar); err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.SuccessWithMsg(c, "更新成功", nil)
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=50"`
}

// ChangePassword 修改密码
func (ctrl *UserController) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "未授权")
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := ctrl.userService.ChangePassword(userID.(uint), req.OldPassword, req.NewPassword); err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.SuccessWithMsg(c, "密码修改成功", nil)
}

// GetUserByID 根据ID获取用户信息
func (ctrl *UserController) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的用户ID")
		return
	}

	user, err := ctrl.userService.GetUserByID(uint(id))
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Success(c, user)
}
