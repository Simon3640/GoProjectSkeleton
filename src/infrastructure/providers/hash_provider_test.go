package providers

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPasswordAndVerifyPassword(t *testing.T) {
	assert := assert.New(t)

	hashProvider := NewHashProvider()

	password := "StrongP@ssw0rd!"

	hashedPassword, err := hashProvider.HashPassword(password)

	assert.Nil(err)
	assert.NotEmpty(hashedPassword)

	// prefix
	assert.True(strings.HasPrefix(hashedPassword, "$argon2id$"))

	ok, err := hashProvider.VerifyPassword(hashedPassword, password)
	assert.Nil(err)
	assert.True(ok)

	wrongPassword := "WrongP@ssw0rd!"
	ok, err = hashProvider.VerifyPassword(hashedPassword, wrongPassword)
	assert.Nil(err)
	assert.False(ok)

}

func TestVerifyPasswordWrongHashFormat(t *testing.T) {
	assert := assert.New(t)

	hashProvider := NewHashProvider()

	wrongHash := "invalid$hash$format"
	ok, err := hashProvider.VerifyPassword(wrongHash, "SomePassword")
	assert.NotNil(err)
	assert.False(ok)
}

func TestHashUniqueness(t *testing.T) {
	assert := assert.New(t)

	hashProvider := NewHashProvider()

	password := "AnotherStr0ngP@ss!"
	hash1, err1 := hashProvider.HashPassword(password)
	hash2, err2 := hashProvider.HashPassword(password)
	assert.Nil(err1)
	assert.Nil(err2)
	assert.NotEqual(hash1, hash2, "Hashes should be unique due to different salts")
}

func TestOneTimeToken(t *testing.T) {
	assert := assert.New(t)

	hashProvider := NewHashProvider()

	token1, hash1, err1 := hashProvider.OneTimeToken()
	token2, hash2, err2 := hashProvider.OneTimeToken()

	assert.NotEmpty(token1)
	assert.NotEmpty(token2)
	assert.NotEmpty(hash1)
	assert.NotEmpty(hash2)
	assert.Nil(err1)
	assert.Nil(err2)
	assert.NotEqual(token1, token2, "One-time tokens should be unique")

	// Validate tokens
	valid1 := hashProvider.ValidateOneTimeToken(hash1, token1)
	valid2 := hashProvider.ValidateOneTimeToken(hash2, token2)
	invalid := hashProvider.ValidateOneTimeToken(hash1, token2)
	assert.True(valid1, "Token 1 should be valid")
	assert.True(valid2, "Token 2 should be valid")
	assert.False(invalid, "Token 2 should not validate against hash 1")

	// Test hashing the token
	hashedToken1 := hashProvider.HashOneTimeToken(token1)
	assert.Equal(hash1, hashedToken1, "Hashed token should match the original hash")
}
