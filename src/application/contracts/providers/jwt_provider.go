package contractsproviders

import (
	"context"
	"time"

	application_errors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
)

type JWTCLaims map[string]interface{}

type IJWTProvider interface {
	GenerateAccessToken(ctx context.Context, subject string, claimsMap JWTCLaims) (string, time.Time, *application_errors.ApplicationError)
	GenerateRefreshToken(ctx context.Context, subject string) (string, time.Time, *application_errors.ApplicationError)
	ParseTokenAndValidate(tokenString string) (JWTCLaims, *application_errors.ApplicationError)
}
