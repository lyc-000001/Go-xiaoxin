package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/xiaoxin/blog-backend/pkg/config"
)

// SaveUploadedFile 保存上传的文件
func SaveUploadedFile(file *multipart.FileHeader) (string, error) {
	cfg := config.GlobalConfig.Upload

	// 检查文件大小
	if file.Size > int64(cfg.MaxSize)*1024*1024 {
		return "", fmt.Errorf("文件大小超过限制，最大允许 %dMB", cfg.MaxSize)
	}

	// 检查文件扩展名
	ext := strings.ToLower(path.Ext(file.Filename))
	if !isAllowedExt(ext, cfg.AllowedExts) {
		return "", fmt.Errorf("不支持的文件类型: %s", ext)
	}

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("打开上传文件失败: %w", err)
	}
	defer src.Close()

	// 生成唯一文件名
	filename := generateUniqueFilename(ext)

	// 创建日期目录
	dateDir := time.Now().Format("2006-01-02")
	savePath := filepath.Join(cfg.SavePath, dateDir)
	if err := os.MkdirAll(savePath, os.ModePerm); err != nil {
		return "", fmt.Errorf("创建目录失败: %w", err)
	}

	// 完整的文件路径
	fullPath := filepath.Join(savePath, filename)

	// 创建目标文件
	dst, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("创建文件失败: %w", err)
	}
	defer dst.Close()

	// 复制文件内容
	if _, err = io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("保存文件失败: %w", err)
	}

	// 返回相对路径
	relativePath := filepath.Join(dateDir, filename)
	return relativePath, nil
}

// isAllowedExt 检查文件扩展名是否允许
func isAllowedExt(ext string, allowedExts []string) bool {
	for _, allowedExt := range allowedExts {
		if ext == allowedExt {
			return true
		}
	}
	return false
}

// generateUniqueFilename 生成唯一文件名
func generateUniqueFilename(ext string) string {
	return uuid.New().String() + ext
}

// DeleteFile 删除文件
func DeleteFile(relativePath string) error {
	cfg := config.GlobalConfig.Upload
	fullPath := filepath.Join(cfg.SavePath, relativePath)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return fmt.Errorf("文件不存在")
	}

	if err := os.Remove(fullPath); err != nil {
		return fmt.Errorf("删除文件失败: %w", err)
	}

	return nil
}

// GetFileURL 获取文件访问URL
func GetFileURL(relativePath string) string {
	if relativePath == "" {
		return ""
	}
	cfg := config.GlobalConfig
	return fmt.Sprintf("http://localhost:%d/uploads/%s", cfg.App.Port, strings.ReplaceAll(relativePath, "\\", "/"))
}
