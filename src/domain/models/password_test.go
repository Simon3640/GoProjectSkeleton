package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasswordModel(t *testing.T) {
	assert := assert.New(t)

	// Test NewPasswordCreate
	passwordCreate := NewPasswordCreate(1, "TestPassword123", nil, true)
	assert.NotNil(passwordCreate)
	assert.Equal(passwordCreate.UserID, uint(1))
	assert.Equal(passwordCreate.Hash, "TestPassword123")
	assert.NotNil(passwordCreate.ExpiresAt)
	assert.True(passwordCreate.IsActive)

	// Test PasswordNoHash IsValidPassword
	validPassword := PasswordCreateNoHash{
		UserID:           1,
		NoHashedPassword: "ValidPass123!",
		IsActive:         true,
	}
	assert.True(validPassword.IsValidPassword())

	invalidPassword := PasswordCreateNoHash{
		UserID:           1,
		NoHashedPassword: "short",
		IsActive:         true,
	}
	assert.False(invalidPassword.IsValidPassword())

	invalidPassword2 := PasswordCreateNoHash{
		UserID:           1,
		NoHashedPassword: "nouppercase123!",
		IsActive:         true,
	}
	assert.False(invalidPassword2.IsValidPassword())

	invalidPassword3 := PasswordCreateNoHash{
		UserID:           1,
		NoHashedPassword: "NOLOWERCASE123!",
		IsActive:         true,
	}
	assert.False(invalidPassword3.IsValidPassword())

	invalidPassword4 := PasswordCreateNoHash{
		UserID:           1,
		NoHashedPassword: "NoNumber!",
		IsActive:         true,
	}
	assert.False(invalidPassword4.IsValidPassword())

	invalidPassword5 := PasswordCreateNoHash{
		UserID:           1,
		NoHashedPassword: "NoSpecial123",
		IsActive:         true,
	}
	assert.False(invalidPassword5.IsValidPassword())

}
