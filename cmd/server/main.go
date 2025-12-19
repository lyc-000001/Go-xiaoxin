package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"github.com/xiaoxin/blog-backend/internal/middleware"
	"github.com/xiaoxin/blog-backend/internal/models"
	"github.com/xiaoxin/blog-backend/internal/routes"
	"github.com/xiaoxin/blog-backend/pkg/config"
	"github.com/xiaoxin/blog-backend/pkg/database"
	pkgjwt "github.com/xiaoxin/blog-backend/pkg/jwt"
	"github.com/xiaoxin/blog-backend/pkg/logger"
	"github.com/xiaoxin/blog-backend/pkg/redis"
)

func main() {
	// 1. 加载配置
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	// 2. 初始化日志系统
	if err := logger.InitLogger(&cfg.Log); err != nil {
		log.Fatalf("初始化日志系统失败: %v", err)
	}
	defer logger.Sync()

	logger.Info("日志系统初始化完成")

	// 3. 初始化数据库
	if err := database.InitDB(&cfg.Database); err != nil {
		logger.Fatalf("初始化数据库失败: %v", err)
	}
	defer database.CloseDB()

	// 自动迁移数据表
	if err := database.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Tag{},
		&models.Article{},
		&models.Comment{},
	); err != nil {
		logger.Fatalf("数据表迁移失败: %v", err)
	}
	logger.Info("数据表迁移完成")

	// 4. 初始化Redis
	if err := redis.InitRedis(&cfg.Redis); err != nil {
		logger.Fatalf("初始化Redis失败: %v", err)
	}
	defer redis.CloseRedis()

	// 5. 初始化JWT
	pkgjwt.InitJWT(cfg.JWT.Secret)
	logger.Info("JWT初始化完成")

	// 6. 设置Gin模式
	gin.SetMode(cfg.App.Mode)

	// 7. 创建路由引擎
	r := gin.New()

	// 8. 使用中间件
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())

	// 9. 设置路由
	routes.SetupRoutes(r)

	// 10. 启动服务
	addr := fmt.Sprintf(":%d", cfg.App.Port)
	logger.Infof("服务启动成功，监听地址: %s", addr)
	logger.Infof("应用名称: %s", cfg.App.Name)
	logger.Infof("应用版本: %s", cfg.App.Version)
	logger.Infof("运行模式: %s", cfg.App.Mode)

	// 11. 优雅关闭
	go func() {
		if err := r.Run(addr); err != nil {
			logger.Fatalf("服务启动失败: %v", err)
		}
	}()

	// 12. 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("正在关闭服务...")
	logger.Info("服务已关闭")
}
