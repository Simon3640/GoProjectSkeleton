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
