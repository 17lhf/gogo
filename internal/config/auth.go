package config

import "time"

// AuthConfig holds JWT and session related configuration.
type AuthConfig struct {
	JWTSecret         string
	SessionTTL        time.Duration
	LockoutThreshold  int
	LockoutDuration   time.Duration
	PasswordCost      int
	PasswordMaxAge    time.Duration
}

func loadAuth() AuthConfig {
	return AuthConfig{
		JWTSecret:        getenv("JWT_SECRET", "gogo-dev-secret-change-in-production"),
		SessionTTL:       8 * time.Hour,
		LockoutThreshold: 5,
		LockoutDuration:  30 * time.Minute,
		PasswordCost:     12,
		PasswordMaxAge:   365 * 24 * time.Hour,
	}
}
