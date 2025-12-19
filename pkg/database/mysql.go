package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/xiaoxin/blog-backend/pkg/config"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB(cfg *config.DatabaseConfig) error {
	dsn := cfg.GetDSN()

	// 设置日志级别
	var logLevel logger.LogLevel
	switch cfg.LogLevel {
	case 1:
		logLevel = logger.Silent
	case 2:
		logLevel = logger.Error
	case 3:
		logLevel = logger.Warn
	case 4:
		logLevel = logger.Info
	default:
		logLevel = logger.Info
	}

	// 创建数据库连接
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	})
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	// 获取底层的 sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %w", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("数据库连接测试失败: %w", err)
	}

	DB = db
	log.Println("数据库连接成功")
	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}

// CloseDB 关闭数据库连接
func CloseDB() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// AutoMigrate 自动迁移数据表
func AutoMigrate(models ...interface{}) error {
	if DB == nil {
		return fmt.Errorf("数据库未初始化")
	}
	return DB.AutoMigrate(models...)
}
