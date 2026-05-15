package config

import "github.com/spf13/viper"

// ServerConfig holds HTTP server parameters.
type ServerConfig struct {
	Port string
}

func loadServer() ServerConfig {
	return ServerConfig{
		Port: viper.GetString("APP_PORT"),
	}
}
