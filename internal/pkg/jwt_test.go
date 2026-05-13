package pkg

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateAndParseToken(t *testing.T) {
	secret := "test-secret"
	roles := []string{"OPERATOR"}

	token, jti, err := GenerateToken(secret, 1, "testuser", roles)
	require.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.NotEmpty(t, jti)

	claims, err := ParseToken(secret, token)
	require.NoError(t, err)
	assert.Equal(t, int64(1), claims.UserID)
	assert.Equal(t, "testuser", claims.Username)
	assert.Equal(t, roles, claims.Roles)
	assert.Equal(t, jti, claims.ID)
	assert.True(t, claims.ExpiresAt.Time.After(time.Now()))
}

func TestParseToken_InvalidSecret(t *testing.T) {
	token, _, err := GenerateToken("secret-a", 1, "test", []string{})
	require.NoError(t, err)

	_, err = ParseToken("secret-b", token)
	assert.Error(t, err)
}

func TestParseToken_InvalidToken(t *testing.T) {
	_, err := ParseToken("secret", "not.a.valid.token")
	assert.Error(t, err)
}
