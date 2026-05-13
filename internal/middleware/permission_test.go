package middleware

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasRole_SuperAdmin(t *testing.T) {
	roles := []string{SuperAdminCode, "OPERATOR"}
	assert.True(t, HasRole(roles, SuperAdminCode))
}

func TestHasRole_NotFound(t *testing.T) {
	roles := []string{"OPERATOR"}
	assert.False(t, HasRole(roles, SuperAdminCode))
}

func TestHasRole_Empty(t *testing.T) {
	assert.False(t, HasRole(nil, SuperAdminCode))
	assert.False(t, HasRole([]string{}, SuperAdminCode))
}

func TestSuperAdminCode_Constant(t *testing.T) {
	assert.Equal(t, "SUPER_ADMIN", SuperAdminCode)
}
