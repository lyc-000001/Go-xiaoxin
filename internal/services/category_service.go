package services

import (
	"errors"

	"gorm.io/gorm"

	"github.com/xiaoxin/blog-backend/internal/models"
	"github.com/xiaoxin/blog-backend/pkg/database"
)

// CategoryService 分类服务
type CategoryService struct{}

// NewCategoryService 创建分类服务实例
func NewCategoryService() *CategoryService {
	return &CategoryService{}
}

// CreateCategory 创建分类
func (s *CategoryService) CreateCategory(name, description string, sort int) (*models.Category, error) {
	db := database.GetDB()

	// 检查分类名是否存在
	var count int64
	if err := db.Model(&models.Category{}).Where("name = ?", name).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("分类名已存在")
	}

	category := &models.Category{
		Name:        name,
		Description: description,
		Sort:        sort,
	}

	if err := db.Create(category).Error; err != nil {
		return nil, err
	}

	return category, nil
}

// GetCategoryByID 根据ID获取分类
func (s *CategoryService) GetCategoryByID(id uint) (*models.Category, error) {
	db := database.GetDB()

	var category models.Category
	if err := db.First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("分类不存在")
		}
		return nil, err
	}

	return &category, nil
}

// GetCategoryList 获取分类列表
func (s *CategoryService) GetCategoryList() ([]models.Category, error) {
	db := database.GetDB()

	var categories []models.Category
	if err := db.Order("sort ASC, id DESC").Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

// UpdateCategory 更新分类
func (s *CategoryService) UpdateCategory(id uint, name, description string, sort int) error {
	db := database.GetDB()

	// 检查分类是否存在
	var category models.Category
	if err := db.First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("分类不存在")
		}
		return err
	}

	// 更新分类
	updates := map[string]interface{}{
		"name":        name,
		"description": description,
		"sort":        sort,
	}

	if err := db.Model(&models.Category{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return err
	}

	return nil
}

// DeleteCategory 删除分类
func (s *CategoryService) DeleteCategory(id uint) error {
	db := database.GetDB()

	// 检查分类是否存在
	var category models.Category
	if err := db.First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("分类不存在")
		}
		return err
	}

	// 检查是否有文章使用该分类
	var count int64
	if err := db.Model(&models.Article{}).Where("category_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("该分类下还有文章，无法删除")
	}

	// 删除分类
	if err := db.Delete(&category).Error; err != nil {
		return err
	}

	return nil
}
