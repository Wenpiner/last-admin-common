package last_redis

type RedisKeyPrefix string

const (
	CasbinChannel RedisKeyPrefix = "casbin:"
	BlacklistToken RedisKeyPrefix = "blacklist:tokens"
)
