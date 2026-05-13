package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("Test1234")
	require.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, "Test1234", hash)
}

func TestCheckPassword(t *testing.T) {
	hash, err := HashPassword("Test1234")
	require.NoError(t, err)

	err = CheckPassword(hash, "Test1234")
	assert.NoError(t, err)

	err = CheckPassword(hash, "WrongPassword1")
	assert.Error(t, err)
}

func TestValidatePasswordStrength(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
		errMsg   string
	}{
		{"valid password", "Test1234", false, ""},
		{"too short", "Te1", true, "密码长度不能少于8位"},
		{"no uppercase", "test1234", true, "密码必须包含大写字母"},
		{"no lowercase", "TEST1234", true, "密码必须包含小写字母"},
		{"no digit", "TestTest", true, "密码必须包含数字"},
		{"exactly 8 chars", "Test1234", false, ""},
		{"special chars ok", "Test@#$%1", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePasswordStrength(tt.password)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
