package contracts

import (
	"context"
	"time"
)

type JWTCLaims map[string]interface{}

type IJWTProvider interface {
	GenerateAccessToken(ctx context.Context, subject string, claimsMap JWTCLaims) (string, time.Time, error)
	GenerateRefreshToken(ctx context.Context, subject string) (string, time.Time, error)
	ParseTokenAndValidate(tokenString string) (JWTCLaims, error)
}
