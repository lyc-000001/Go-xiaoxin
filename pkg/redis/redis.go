package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/xiaoxin/blog-backend/pkg/config"
)

var Client *redis.Client
var Ctx = context.Background()

// InitRedis 初始化Redis连接
func InitRedis(cfg *config.RedisConfig) error {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.GetRedisAddr(),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	// 测试连接
	_, err := client.Ping(Ctx).Result()
	if err != nil {
		return fmt.Errorf("Redis连接失败: %w", err)
	}

	Client = client
	log.Println("Redis连接成功")
	return nil
}

// GetClient 获取Redis客户端
func GetClient() *redis.Client {
	return Client
}

// CloseRedis 关闭Redis连接
func CloseRedis() error {
	if Client != nil {
		return Client.Close()
	}
	return nil
}

// Set 设置键值对
func Set(key string, value interface{}, expiration time.Duration) error {
	return Client.Set(Ctx, key, value, expiration).Err()
}

// Get 获取值
func Get(key string) (string, error) {
	return Client.Get(Ctx, key).Result()
}

// Delete 删除键
func Delete(keys ...string) error {
	return Client.Del(Ctx, keys...).Err()
}

// Exists 检查键是否存在
func Exists(keys ...string) (int64, error) {
	return Client.Exists(Ctx, keys...).Result()
}

// Expire 设置过期时间
func Expire(key string, expiration time.Duration) error {
	return Client.Expire(Ctx, key, expiration).Err()
}

// GetTTL 获取键的剩余生存时间
func GetTTL(key string) (time.Duration, error) {
	return Client.TTL(Ctx, key).Result()
}

// Incr 自增
func Incr(key string) (int64, error) {
	return Client.Incr(Ctx, key).Result()
}

// Decr 自减
func Decr(key string) (int64, error) {
	return Client.Decr(Ctx, key).Result()
}

// SetNX 设置键值对（仅当键不存在时）
func SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	return Client.SetNX(Ctx, key, value, expiration).Result()
}

// HSet 设置哈希字段
func HSet(key string, values ...interface{}) error {
	return Client.HSet(Ctx, key, values...).Err()
}

// HGet 获取哈希字段值
func HGet(key, field string) (string, error) {
	return Client.HGet(Ctx, key, field).Result()
}

// HGetAll 获取所有哈希字段和值
func HGetAll(key string) (map[string]string, error) {
	return Client.HGetAll(Ctx, key).Result()
}

// HDel 删除哈希字段
func HDel(key string, fields ...string) error {
	return Client.HDel(Ctx, key, fields...).Err()
}
