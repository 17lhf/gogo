package config

// RedisConfig holds connection parameters for Redis.
type RedisConfig struct {
	Addr     string
	Password string
}

func loadRedis() RedisConfig {
	return RedisConfig{
		Addr:     getenv("REDIS_ADDR", "localhost:6379"),
		Password: getenv("REDIS_PASSWORD", ""),
	}
}
