package controllers

import (
	"github.com/gin-gonic/gin"

	"github.com/xiaoxin/blog-backend/internal/utils"
)

// UploadController 文件上传控制器
type UploadController struct{}

// NewUploadController 创建文件上传控制器实例
func NewUploadController() *UploadController {
	return &UploadController{}
}

// UploadFile 上传文件
func (ctrl *UploadController) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		utils.BadRequest(c, "请选择要上传的文件")
		return
	}

	// 保存文件
	relativePath, err := utils.SaveUploadedFile(file)
	if err != nil {
		utils.Error(c, err.Error())
		return
	}

	// 获取文件访问URL
	fileURL := utils.GetFileURL(relativePath)

	utils.Success(c, gin.H{
		"path": relativePath,
		"url":  fileURL,
	})
}
