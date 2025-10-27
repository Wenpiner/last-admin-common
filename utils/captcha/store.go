package captcha

import (
	"context"
	"fmt"
	"time"

	"github.com/mojocn/base64Captcha"
	"github.com/redis/go-redis/v9"
)

// Store 验证码存储接口
type Store interface {
	base64Captcha.Store
	// Close 关闭存储连接
	Close() error
}

// RedisStore Redis存储实现
type RedisStore struct {
	client    *redis.Client
	keyPrefix string
	expire    time.Duration
	ctx       context.Context
}

// NewRedisStore 创建Redis存储
func NewRedisStore(config RedisConfig, keyPrefix string, expire time.Duration) (*RedisStore, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
		PoolSize: config.PoolSize,
	})

	// 测试连接
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	return &RedisStore{
		client:    client,
		keyPrefix: keyPrefix,
		expire:    expire,
		ctx:       ctx,
	}, nil
}

// Set 存储验证码
func (r *RedisStore) Set(id string, value string) error {
	key := r.getKey(id)
	return r.client.Set(r.ctx, key, value, r.expire).Err()
}

// Get 获取验证码
func (r *RedisStore) Get(id string, clear bool) string {
	key := r.getKey(id)
	
	val, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		return ""
	}
	
	if clear {
		r.client.Del(r.ctx, key)
	}
	
	return val
}

// Verify 验证验证码
func (r *RedisStore) Verify(id, answer string, clear bool) bool {
	key := r.getKey(id)
	
	val, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		return false
	}
	
	if clear {
		r.client.Del(r.ctx, key)
	}
	
	return val == answer
}

// Close 关闭Redis连接
func (r *RedisStore) Close() error {
	return r.client.Close()
}

// Clear 清除验证码
func (r *RedisStore) Clear(id string) {
	r.client.Del(r.ctx, r.getKey(id))
}

// getKey 获取完整的Redis key
func (r *RedisStore) getKey(id string) string {
	return r.keyPrefix + id
}

// MemoryStore 内存存储实现
type MemoryStore struct {
	store base64Captcha.Store
}

// NewMemoryStore 创建内存存储
func NewMemoryStore(collectNum int, expiration time.Duration) *MemoryStore {
	return &MemoryStore{
		store: base64Captcha.NewMemoryStore(collectNum, expiration),
	}
}

// Set 存储验证码
func (m *MemoryStore) Set(id string, value string) error {
	return m.store.Set(id, value)
}

// Get 获取验证码
func (m *MemoryStore) Get(id string, clear bool) string {
	return m.store.Get(id, clear)
}

// Verify 验证验证码
func (m *MemoryStore) Verify(id, answer string, clear bool) bool {
	return m.store.Verify(id, answer, clear)
}

// Clear 清除验证码
func (m *MemoryStore) Clear(id string) {
	m.store.Get(id, true)
}

// Close 关闭内存存储（空实现）
func (m *MemoryStore) Close() error {
	return nil
}

// NewStore 根据配置创建存储实例
func NewStore(config StoreConfig) (Store, error) {
	switch config.Type {
	case StoreTypeRedis:
		return NewRedisStore(config.Redis, config.KeyPrefix, config.Expire)
	case StoreTypeMemory:
		return NewMemoryStore(1000, config.Expire), nil
	default:
		return nil, fmt.Errorf("unsupported store type: %s", config.Type)
	}
}
