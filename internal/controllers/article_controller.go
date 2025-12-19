package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/xiaoxin/blog-backend/internal/models"
	"github.com/xiaoxin/blog-backend/internal/services"
	"github.com/xiaoxin/blog-backend/internal/utils"
)

// ArticleController 文章控制器
type ArticleController struct {
	articleService *services.ArticleService
}

// NewArticleController 创建文章控制器实例
func NewArticleController() *ArticleController {
	return &ArticleController{
		articleService: services.NewArticleService(),
	}
}

// CreateArticleRequest 创建文章请求
type CreateArticleRequest struct {
	Title       string `json:"title" binding:"required,max=200"`
	Description string `json:"description" binding:"max=500"`
	Content     string `json:"content" binding:"required"`
	Cover       string `json:"cover" binding:"max=255"`
	CategoryID  uint   `json:"category_id"`
	TagIDs      []uint `json:"tag_ids"`
	Status      int    `json:"status" binding:"oneof=0 1"`
	IsTop       bool   `json:"is_top"`
}

// CreateArticle 创建文章
func (ctrl *ArticleController) CreateArticle(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "未授权")
		return
	}

	var req CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 构建标签
	var tags []models.Tag
	for _, tagID := range req.TagIDs {
		tags = append(tags, models.Tag{BaseModel: models.BaseModel{ID: tagID}})
	}

	article := &models.Article{
		Title:       req.Title,
		Description: req.Description,
		Content:     req.Content,
		Cover:       req.Cover,
		AuthorID:    userID.(uint),
		CategoryID:  req.CategoryID,
		Tags:        tags,
		Status:      req.Status,
		IsTop:       req.IsTop,
	}

	if err := ctrl.articleService.CreateArticle(article); err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.SuccessWithMsg(c, "创建成功", gin.H{"id": article.ID})
}

// GetArticle 获取文章详情
func (ctrl *ArticleController) GetArticle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的文章ID")
		return
	}

	article, err := ctrl.articleService.GetArticleByID(uint(id))
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	// 增加浏览量
	_ = ctrl.articleService.IncrementViewCount(uint(id))

	utils.Success(c, article)
}

// GetArticleList 获取文章列表
func (ctrl *ArticleController) GetArticleList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	var status *int
	if statusStr := c.Query("status"); statusStr != "" {
		s, _ := strconv.Atoi(statusStr)
		status = &s
	}

	var categoryID *uint
	if categoryIDStr := c.Query("category_id"); categoryIDStr != "" {
		cid, _ := strconv.ParseUint(categoryIDStr, 10, 32)
		catID := uint(cid)
		categoryID = &catID
	}

	articles, total, err := ctrl.articleService.GetArticleList(page, pageSize, status, categoryID)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.PageSuccess(c, articles, total, page, pageSize)
}

// UpdateArticleRequest 更新文章请求
type UpdateArticleRequest struct {
	Title       string `json:"title" binding:"max=200"`
	Description string `json:"description" binding:"max=500"`
	Content     string `json:"content"`
	Cover       string `json:"cover" binding:"max=255"`
	CategoryID  uint   `json:"category_id"`
	TagIDs      []uint `json:"tag_ids"`
	Status      int    `json:"status" binding:"oneof=0 1"`
	IsTop       bool   `json:"is_top"`
}

// UpdateArticle 更新文章
func (ctrl *ArticleController) UpdateArticle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的文章ID")
		return
	}

	var req UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 构建标签
	var tags []models.Tag
	for _, tagID := range req.TagIDs {
		tags = append(tags, models.Tag{BaseModel: models.BaseModel{ID: tagID}})
	}

	article := &models.Article{
		Title:       req.Title,
		Description: req.Description,
		Content:     req.Content,
		Cover:       req.Cover,
		CategoryID:  req.CategoryID,
		Tags:        tags,
		Status:      req.Status,
		IsTop:       req.IsTop,
	}

	if err := ctrl.articleService.UpdateArticle(uint(id), article); err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.SuccessWithMsg(c, "更新成功", nil)
}

// DeleteArticle 删除文章
func (ctrl *ArticleController) DeleteArticle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的文章ID")
		return
	}

	if err := ctrl.articleService.DeleteArticle(uint(id)); err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.SuccessWithMsg(c, "删除成功", nil)
}

// LikeArticle 点赞文章
func (ctrl *ArticleController) LikeArticle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的文章ID")
		return
	}

	if err := ctrl.articleService.IncrementLikeCount(uint(id)); err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.SuccessWithMsg(c, "点赞成功", nil)
}
