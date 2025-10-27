package redis

type RedisKeyPrefix string

const (
	CasbinChannel RedisKeyPrefix = "casbin:"
)
