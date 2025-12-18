package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/xiaoxin/blog-backend/internal/services"
	"github.com/xiaoxin/blog-backend/internal/utils"
)

// CategoryController 分类控制器
type CategoryController struct {
	categoryService *services.CategoryService
}

// NewCategoryController 创建分类控制器实例
func NewCategoryController() *CategoryController {
	return &CategoryController{
		categoryService: services.NewCategoryService(),
	}
}

// CreateCategoryRequest 创建分类请求
type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required,max=50"`
	Description string `json:"description" binding:"max=255"`
	Sort        int    `json:"sort"`
}

// CreateCategory 创建分类
func (ctrl *CategoryController) CreateCategory(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	category, err := ctrl.categoryService.CreateCategory(req.Name, req.Description, req.Sort)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.SuccessWithMsg(c, "创建成功", category)
}

// GetCategory 获取分类详情
func (ctrl *CategoryController) GetCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的分类ID")
		return
	}

	category, err := ctrl.categoryService.GetCategoryByID(uint(id))
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Success(c, category)
}

// GetCategoryList 获取分类列表
func (ctrl *CategoryController) GetCategoryList(c *gin.Context) {
	categories, err := ctrl.categoryService.GetCategoryList()
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.Success(c, categories)
}

// UpdateCategoryRequest 更新分类请求
type UpdateCategoryRequest struct {
	Name        string `json:"name" binding:"required,max=50"`
	Description string `json:"description" binding:"max=255"`
	Sort        int    `json:"sort"`
}

// UpdateCategory 更新分类
func (ctrl *CategoryController) UpdateCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的分类ID")
		return
	}

	var req UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := ctrl.categoryService.UpdateCategory(uint(id), req.Name, req.Description, req.Sort); err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.SuccessWithMsg(c, "更新成功", nil)
}

// DeleteCategory 删除分类
func (ctrl *CategoryController) DeleteCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequest(c, "无效的分类ID")
		return
	}

	if err := ctrl.categoryService.DeleteCategory(uint(id)); err != nil {
		utils.Error(c, err.Error())
		return
	}

	utils.SuccessWithMsg(c, "删除成功", nil)
}
