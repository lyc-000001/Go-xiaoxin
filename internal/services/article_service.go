package services

import (
	"errors"

	"gorm.io/gorm"

	"github.com/xiaoxin/blog-backend/internal/models"
	"github.com/xiaoxin/blog-backend/pkg/database"
)

// ArticleService 文章服务
type ArticleService struct{}

// NewArticleService 创建文章服务实例
func NewArticleService() *ArticleService {
	return &ArticleService{}
}

// CreateArticle 创建文章
func (s *ArticleService) CreateArticle(article *models.Article) error {
	db := database.GetDB()
	return db.Create(article).Error
}

// GetArticleByID 根据ID获取文章
func (s *ArticleService) GetArticleByID(id uint) (*models.Article, error) {
	db := database.GetDB()

	var article models.Article
	if err := db.Preload("Author").Preload("Category").Preload("Tags").First(&article, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("文章不存在")
		}
		return nil, err
	}

	return &article, nil
}

// GetArticleList 获取文章列表
func (s *ArticleService) GetArticleList(page, pageSize int, status *int, categoryID *uint) ([]models.Article, int64, error) {
	db := database.GetDB()

	var articles []models.Article
	var total int64

	query := db.Model(&models.Article{})

	// 筛选条件
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if categoryID != nil {
		query = query.Where("category_id = ?", *categoryID)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Preload("Author").Preload("Category").Preload("Tags").
		Order("is_top DESC, created_at DESC").
		Limit(pageSize).Offset(offset).
		Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// UpdateArticle 更新文章
func (s *ArticleService) UpdateArticle(id uint, article *models.Article) error {
	db := database.GetDB()

	// 检查文章是否存在
	var existingArticle models.Article
	if err := db.First(&existingArticle, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("文章不存在")
		}
		return err
	}

	// 更新文章
	article.ID = id
	if err := db.Model(&models.Article{}).Where("id = ?", id).Updates(article).Error; err != nil {
		return err
	}

	// 更新标签关联
	if len(article.Tags) > 0 {
		if err := db.Model(&existingArticle).Association("Tags").Replace(article.Tags); err != nil {
			return err
		}
	}

	return nil
}

// DeleteArticle 删除文章
func (s *ArticleService) DeleteArticle(id uint) error {
	db := database.GetDB()

	// 检查文章是否存在
	var article models.Article
	if err := db.First(&article, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("文章不存在")
		}
		return err
	}

	// 删除文章
	if err := db.Delete(&article).Error; err != nil {
		return err
	}

	return nil
}

// IncrementViewCount 增加浏览量
func (s *ArticleService) IncrementViewCount(id uint) error {
	db := database.GetDB()
	return db.Model(&models.Article{}).Where("id = ?", id).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

// IncrementLikeCount 增加点赞数
func (s *ArticleService) IncrementLikeCount(id uint) error {
	db := database.GetDB()
	return db.Model(&models.Article{}).Where("id = ?", id).UpdateColumn("like_count", gorm.Expr("like_count + ?", 1)).Error
}
