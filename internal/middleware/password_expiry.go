package middleware

import (
	"time"

	"github.com/gin-gonic/gin"

	"gogo/internal/config"
	"gogo/internal/i18n"
	"gogo/internal/pkg"
	"gogo/internal/repository"
)

// PasswordExpiry returns a middleware that checks password expiry.
// It allows GET /api/v1/auth/me and PUT /api/v1/auth/password regardless of expiry.
func PasswordExpiry(userRepo repository.UserRepository, cfg config.AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := GetUserID(c)
		if userID == 0 {
			c.Next()
			return
		}

		path := c.FullPath()
		// Allow password change and me endpoint unconditionally
		if path == "/api/v1/auth/password" || path == "/api/v1/auth/me" {
			c.Next()
			return
		}

		user, err := userRepo.GetByID(c.Request.Context(), userID)
		if err != nil || user == nil {
			c.Next()
			return
		}

		// Check must_change_password flag
		if user.MustChangePassword {
			pkg.Error(c, 403, pkg.CodeMustChangePassword, i18n.Localize(c, i18n.MsgMustChangePassword))
			return
		}

		// Check password age
		if time.Since(user.PasswordUpdatedAt) > cfg.PasswordMaxAge {
			pkg.Error(c, 403, pkg.CodePasswordExpired, i18n.Localize(c, i18n.MsgPasswordExpired))
			return
		}

		c.Next()
	}
}
