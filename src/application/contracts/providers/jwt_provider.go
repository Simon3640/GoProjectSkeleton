package contractsproviders

import (
	"context"
	application_errors "gormgoskeleton/src/application/shared/errors"
	"time"
)

type JWTCLaims map[string]interface{}

type IJWTProvider interface {
	GenerateAccessToken(ctx context.Context, subject string, claimsMap JWTCLaims) (string, time.Time, *application_errors.ApplicationError)
	GenerateRefreshToken(ctx context.Context, subject string) (string, time.Time, *application_errors.ApplicationError)
	ParseTokenAndValidate(tokenString string) (JWTCLaims, *application_errors.ApplicationError)
}
