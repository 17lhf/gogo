package config

import "os"

// Config holds all application configuration loaded from environment variables.
type Config struct {
	Postgres PostgresConfig
	Redis    RedisConfig
	Server   ServerConfig
}

// Load reads all environment variables and returns a populated Config.
func Load() Config {
	return Config{
		Postgres: loadPostgres(),
		Redis:    loadRedis(),
		Server:   loadServer(),
	}
}

// getenv returns the env variable value, or fallback if not set / empty.
func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
