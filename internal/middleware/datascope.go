package middleware

import (
	"github.com/gin-gonic/gin"

	"gogo/internal/repository"
)

const ContextKeyStoreIDs = "store_ids"

// DataScope returns a middleware that injects the user's accessible store IDs
// into the request context. SUPER_ADMIN users get nil (access to all stores).
func DataScope(userRepo repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := GetUserID(c)
		roles := GetRoles(c)

		// SUPER_ADMIN sees all stores
		if HasRole(roles, SuperAdminCode) {
			c.Set(ContextKeyStoreIDs, []int64{})
			c.Next()
			return
		}

		storeIDs, err := userRepo.GetStoreIDs(c.Request.Context(), userID)
		if err != nil {
			// On error, provide empty scope (no access)
			c.Set(ContextKeyStoreIDs, []int64{})
			c.Next()
			return
		}

		c.Set(ContextKeyStoreIDs, storeIDs)
		c.Next()
	}
}

// GetStoreIDs extracts store IDs from the Gin context.
func GetStoreIDs(c *gin.Context) []int64 {
	v, _ := c.Get(ContextKeyStoreIDs)
	if ids, ok := v.([]int64); ok {
		return ids
	}
	return nil
}
