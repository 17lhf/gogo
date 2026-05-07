package config

// ServerConfig holds HTTP server parameters.
type ServerConfig struct {
	Port string
}

func loadServer() ServerConfig {
	return ServerConfig{
		Port: getenv("APP_PORT", "8080"),
	}
}
