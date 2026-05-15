package config

import "github.com/spf13/viper"

// RedisConfig holds connection parameters for Redis.
type RedisConfig struct {
	Addr     string
	Password string
}

func loadRedis() RedisConfig {
	return RedisConfig{
		Addr:     viper.GetString("REDIS_ADDR"),
		Password: viper.GetString("REDIS_PASSWORD"),
	}
}
