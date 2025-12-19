// Package authcontracts contains the contracts for the auth module.
package authcontracts

import (
	"context"
	"time"

	applicationerrors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
)

// JWTCLaims is the claims for the JWT token
type JWTCLaims map[string]interface{}

// IJWTProvider is the interface for the JWT provider
type IJWTProvider interface {
	GenerateAccessToken(ctx context.Context, subject string, claimsMap JWTCLaims) (string, time.Time, *applicationerrors.ApplicationError)
	GenerateRefreshToken(ctx context.Context, subject string) (string, time.Time, *applicationerrors.ApplicationError)
	ParseTokenAndValidate(tokenString string) (JWTCLaims, *applicationerrors.ApplicationError)
}
