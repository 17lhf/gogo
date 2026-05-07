package config

import "fmt"

// PostgresConfig holds connection parameters for PostgreSQL.
type PostgresConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	DB       string
}

func loadPostgres() PostgresConfig {
	return PostgresConfig{
		User:     getenv("POSTGRES_USER", "gogo"),
		Password: getenv("POSTGRES_PASSWORD", "gogo123"),
		Host:     getenv("POSTGRES_HOST", "localhost"),
		Port:     getenv("POSTGRES_PORT", "5432"),
		DB:       getenv("POSTGRES_DB", "gogo_dev"),
	}
}

// DSN returns a postgres connection string.
func (c PostgresConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, c.DB,
	)
}
