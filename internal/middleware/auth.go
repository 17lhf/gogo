package middleware

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"gogo/internal/cache"
	"gogo/internal/config"
	"gogo/internal/pkg"
)

const (
	ContextKeyUserID    = "user_id"
	ContextKeyUsername  = "username"
	ContextKeyRoles     = "roles"
	ContextKeySessionID = "jti"
)

// Auth returns a middleware that validates JWT and checks Redis session.
func Auth(sessionCache *cache.SessionCache, cfg config.AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			pkg.Error(c, 401, pkg.CodeUnauthorized, "未提供认证信息")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			pkg.Error(c, 401, pkg.CodeUnauthorized, "认证格式错误")
			return
		}

		claims, err := pkg.ParseToken(cfg.JWTSecret, parts[1])
		if err != nil {
			pkg.Error(c, 401, pkg.CodeTokenExpired, "Token已过期或无效")
			return
		}

		// Validate session exists in Redis
		stored, err := sessionCache.Get(c.Request.Context(), claims.UserID, claims.ID)
		if err != nil || stored == nil {
			pkg.Error(c, 401, pkg.CodeSessionNotFound, "会话已失效，请重新登录")
			return
		}

		// Inject user info into context
		c.Set(ContextKeyUserID, claims.UserID)
		c.Set(ContextKeyUsername, claims.Username)
		c.Set(ContextKeyRoles, claims.Roles)
		c.Set(ContextKeySessionID, claims.ID)

		c.Next()
	}
}

// GetUserID extracts the user ID from the Gin context.
func GetUserID(c *gin.Context) int64 {
	v, _ := c.Get(ContextKeyUserID)
	if id, ok := v.(int64); ok {
		return id
	}
	// Handle float64 from JSON unmarshaling in tests
	if f, ok := v.(float64); ok {
		return int64(f)
	}
	return 0
}

// GetUsername extracts the username from the Gin context.
func GetUsername(c *gin.Context) string {
	v, _ := c.Get(ContextKeyUsername)
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

// GetRoles extracts the roles from the Gin context.
func GetRoles(c *gin.Context) []string {
	v, _ := c.Get(ContextKeyRoles)
	if roles, ok := v.([]string); ok {
		return roles
	}
	return nil
}

// GetSessionID extracts the JWT ID from the Gin context.
func GetSessionID(c *gin.Context) string {
	v, _ := c.Get(ContextKeySessionID)
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

// HasRole checks if the roles contain the given code.
func HasRole(roles []string, code string) bool {
	for _, r := range roles {
		if r == code {
			return true
		}
	}
	return false
}

// GetInt64Param extracts an int64 path parameter.
func GetInt64Param(c *gin.Context, name string) (int64, error) {
	return strconv.ParseInt(c.Param(name), 10, 64)
}
