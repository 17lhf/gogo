package config

import (
	"fmt"

	"github.com/spf13/viper"
)

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
		User:     viper.GetString("POSTGRES_USER"),
		Password: viper.GetString("POSTGRES_PASSWORD"),
		Host:     viper.GetString("POSTGRES_HOST"),
		Port:     viper.GetString("POSTGRES_PORT"),
		DB:       viper.GetString("POSTGRES_DB"),
	}
}

// DSN returns a postgres connection string.
func (c PostgresConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, c.DB,
	)
}
