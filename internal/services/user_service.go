package services

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/xiaoxin/blog-backend/internal/models"
	"github.com/xiaoxin/blog-backend/pkg/database"
	pkgjwt "github.com/xiaoxin/blog-backend/pkg/jwt"
)

// UserService 用户服务
type UserService struct{}

// NewUserService 创建用户服务实例
func NewUserService() *UserService {
	return &UserService{}
}

// Register 用户注册
func (s *UserService) Register(username, password, email, nickname string) (*models.User, error) {
	db := database.GetDB()

	// 检查用户名是否存在
	var count int64
	if err := db.Model(&models.User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否存在
	if email != "" {
		if err := db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
			return nil, err
		}
		if count > 0 {
			return nil, errors.New("邮箱已被使用")
		}
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 创建用户
	user := &models.User{
		Username: username,
		Password: string(hashedPassword),
		Email:    email,
		Nickname: nickname,
		Role:     "user",
		Status:   1,
	}

	if err := db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// Login 用户登录
func (s *UserService) Login(username, password string) (string, string, error) {
	db := database.GetDB()

	// 查找用户
	var user models.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", errors.New("用户名或密码错误")
		}
		return "", "", err
	}

	// 检查用户状态
	if user.Status != 1 {
		return "", "", errors.New("用户已被禁用")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", errors.New("用户名或密码错误")
	}

	// 生成JWT令牌
	token, err := pkgjwt.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return "", "", err
	}

	// 生成刷新令牌
	refreshToken, err := pkgjwt.GenerateRefreshToken(user.ID, user.Username, user.Role)
	if err != nil {
		return "", "", err
	}

	return token, refreshToken, nil
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	db := database.GetDB()

	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	return &user, nil
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(id uint, nickname, email, avatar string) error {
	db := database.GetDB()

	updates := make(map[string]interface{})
	if nickname != "" {
		updates["nickname"] = nickname
	}
	if email != "" {
		updates["email"] = email
	}
	if avatar != "" {
		updates["avatar"] = avatar
	}

	if err := db.Model(&models.User{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return err
	}

	return nil
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(id uint, oldPassword, newPassword string) error {
	db := database.GetDB()

	// 获取用户
	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		return err
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.New("原密码错误")
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 更新密码
	if err := db.Model(&models.User{}).Where("id = ?", id).Update("password", string(hashedPassword)).Error; err != nil {
		return err
	}

	return nil
}
