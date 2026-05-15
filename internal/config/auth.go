package config

import (
	"time"

	"github.com/spf13/viper"
)

// AuthConfig holds JWT and session related configuration.
type AuthConfig struct {
	JWTSecret        string
	SessionTTL       time.Duration
	LockoutThreshold int
	LockoutDuration  time.Duration
	PasswordCost     int
	PasswordMaxAge   time.Duration
}

func loadAuth() AuthConfig {
	return AuthConfig{
		JWTSecret:        viper.GetString("JWT_SECRET"),
		SessionTTL:       getDuration("AUTH_SESSION_TTL", 8*time.Hour),
		LockoutThreshold: getInt("AUTH_LOCKOUT_THRESHOLD", 5),
		LockoutDuration:  getDuration("AUTH_LOCKOUT_DURATION", 30*time.Minute),
		PasswordCost:     getInt("AUTH_PASSWORD_COST", 12),
		PasswordMaxAge:   getDuration("AUTH_PASSWORD_MAX_AGE", 365*24*time.Hour),
	}
}
