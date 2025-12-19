package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/xiaoxin/blog-backend/pkg/config"
)

var Logger *zap.Logger
var SugaredLogger *zap.SugaredLogger

// InitLogger 初始化日志系统
func InitLogger(cfg *config.LogConfig) error {
	// 设置日志级别
	var level zapcore.Level
	switch cfg.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	// 编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 日志文件轮转配置
	lumberJackLogger := &lumberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}

	// 创建核心
	core := zapcore.NewTee(
		// 文件输出
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(lumberJackLogger),
			level,
		),
		// 控制台输出
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			level,
		),
	)

	// 创建Logger
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	Logger = logger
	SugaredLogger = logger.Sugar()

	fmt.Println("日志系统初始化成功")
	return nil
}

// GetLogger 获取Logger实例
func GetLogger() *zap.Logger {
	return Logger
}

// GetSugaredLogger 获取SugaredLogger实例
func GetSugaredLogger() *zap.SugaredLogger {
	return SugaredLogger
}

// Sync 同步日志
func Sync() {
	if Logger != nil {
		_ = Logger.Sync()
	}
}

// Debug 调试日志
func Debug(msg string, fields ...zap.Field) {
	Logger.Debug(msg, fields...)
}

// Info 信息日志
func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

// Warn 警告日志
func Warn(msg string, fields ...zap.Field) {
	Logger.Warn(msg, fields...)
}

// Error 错误日志
func Error(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}

// Fatal 致命错误日志
func Fatal(msg string, fields ...zap.Field) {
	Logger.Fatal(msg, fields...)
}

// Debugf 格式化调试日志
func Debugf(template string, args ...interface{}) {
	SugaredLogger.Debugf(template, args...)
}

// Infof 格式化信息日志
func Infof(template string, args ...interface{}) {
	SugaredLogger.Infof(template, args...)
}

// Warnf 格式化警告日志
func Warnf(template string, args ...interface{}) {
	SugaredLogger.Warnf(template, args...)
}

// Errorf 格式化错误日志
func Errorf(template string, args ...interface{}) {
	SugaredLogger.Errorf(template, args...)
}

// Fatalf 格式化致命错误日志
func Fatalf(template string, args ...interface{}) {
	SugaredLogger.Fatalf(template, args...)
}
