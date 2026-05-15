package config

import (
	"time"

	"github.com/spf13/viper"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	Postgres PostgresConfig
	Redis    RedisConfig
	Server   ServerConfig
	Auth     AuthConfig
}

func init() {
	viper.AutomaticEnv()
}

// Load reads all configuration via viper and returns a populated Config.
func Load() Config {
	return Config{
		Postgres: loadPostgres(),
		Redis:    loadRedis(),
		Server:   loadServer(),
		Auth:     loadAuth(),
	}
}

func getDuration(key string, fallback time.Duration) time.Duration {
	v := viper.GetString(key)
	if v == "" {
		return fallback
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return fallback
	}
	return d
}

func getInt(key string, fallback int) int {
	if !viper.IsSet(key) {
		return fallback
	}
	return viper.GetInt(key)
}
