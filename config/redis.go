package config

import (
	"context"
	"crypto/tls"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
)

type RedisConfig struct {
	Host     string `json:""`               // Redis地址
	Password string `json:""`               // Redis密码
	DB       int    `json:",default=0""`    // Redis数据库
	PoolSize int    `json:",default=10"`    // 连接池大小
	Master   string `json:",optional"`    // Redis主节点
	Tls      bool   `json:",default=false"` // 是否启用TLS
}

func (c RedisConfig) Validate() error {
	return nil
}

func (c RedisConfig) NewUniversalRedis() (redis.UniversalClient, error) {
	err := c.Validate()
	if err != nil {
		return nil, err
	}
	opt := &redis.UniversalOptions{
		Addrs:    strings.Split(c.Host, ","),
		Password: c.Password,
		DB:       c.DB,
		PoolSize: c.PoolSize,
	}
	if c.Master != "" {
		opt.MasterName = c.Master
	}
	if c.Tls {
		opt.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
	}
	rds := redis.NewUniversalClient(opt)

	// 执行3s ping 检查是否连接成功
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := rds.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return rds, nil
}

// NewMustUniversalRedis returns a new redis client with must.
func (c RedisConfig) NewMustUniversalRedis() redis.UniversalClient {
	rds, err := c.NewUniversalRedis()
	logx.Must(err)
	return rds
}
