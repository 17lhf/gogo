package pkg

import (
	"errors"
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
		errCheck func(error) bool
	}{
		{"valid password", "Test1234", false, nil},
		{"too short", "Te1", true, func(e error) bool { return errors.Is(e, ErrPasswordTooShort) }},
		{"no uppercase", "test1234", true, func(e error) bool { return errors.Is(e, ErrPasswordNoUpper) }},
		{"no lowercase", "TEST1234", true, func(e error) bool { return errors.Is(e, ErrPasswordNoLower) }},
		{"no digit", "TestTest", true, func(e error) bool { return errors.Is(e, ErrPasswordNoDigit) }},
		{"exactly 8 chars", "Test1234", false, nil},
		{"special chars ok", "Test@#$%1", false, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePasswordStrength(tt.password)
			if tt.wantErr {
				require.Error(t, err)
				assert.True(t, tt.errCheck(err), "unexpected error: %v", err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
