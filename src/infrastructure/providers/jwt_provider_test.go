package providers

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJWTProvider_GenerateAccessToken(t *testing.T) {
	assert := assert.New(t)

	jwtProvider := NewJWTProvider()

	jwtProvider.Setup(
		"secret",
		"test-issuer",
		"test-audience",
		3600,
		86400,
		60,
	)

	ctx := context.Background()
	subject := "test-subject"
	claimsMap := map[string]interface{}{
		"role": "admin",
	}

	token, exp, err := jwtProvider.GenerateAccessToken(ctx, subject, claimsMap)

	assert.NoError(err)
	assert.NotEmpty(token)
	assert.True(exp.After(time.Now()))
}

func TestJWTProvider_GenerateRefreshToken(t *testing.T) {
	assert := assert.New(t)

	jwtProvider := NewJWTProvider()
	jwtProvider.Setup(
		"secret",
		"test-issuer",
		"test-audience",
		3600,
		86400,
		60,
	)
	ctx := context.Background()
	subject := "test-subject"
	token, exp, err := jwtProvider.GenerateRefreshToken(ctx, subject)

	assert.NoError(err)
	assert.NotEmpty(token)
	assert.True(exp.After(time.Now()))
	assert.True(time.Until(exp) <= jwtProvider.config.RefreshTTL+jwtProvider.config.ClockSkew)
}

func TestJWTProvider_ParseTokenAndValidate(t *testing.T) {
	assert := assert.New(t)

	jwtProvider := NewJWTProvider()
	jwtProvider.Setup(
		"secret",
		"test-issuer",
		"test-audience",
		3600,
		86400,
		60,
	)

	ctx := context.Background()
	subject := "test-subject"
	claimsMap := map[string]interface{}{
		"role": "admin",
	}
	token, _, err := jwtProvider.GenerateAccessToken(ctx, subject, claimsMap)

	assert.NoError(err)
	assert.NotEmpty(token)
	parsedClaims, err := jwtProvider.ParseTokenAndValidate(token)
	assert.NoError(err)
	assert.NotNil(parsedClaims)
	assert.Equal(subject, parsedClaims["sub"])
	assert.Equal("admin", parsedClaims["role"])
}
