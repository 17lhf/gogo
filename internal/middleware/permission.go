package middleware

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"

	"gogo/internal/pkg"
)

const SuperAdminCode = "SUPER_ADMIN"

// Permission returns a middleware that enforces RBAC via Casbin.
// Users with SUPER_ADMIN role bypass all permission checks.
func Permission(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles := GetRoles(c)

		// SUPER_ADMIN bypasses Casbin
		if HasRole(roles, SuperAdminCode) {
			c.Next()
			return
		}

		// Enforce for normal roles
		path := c.FullPath()
		method := c.Request.Method

		for _, role := range roles {
			ok, err := enforcer.Enforce(role, path, method)
			if err == nil && ok {
				c.Next()
				return
			}
		}

		pkg.Error(c, 403, pkg.CodeForbidden, "权限不足")
	}
}
