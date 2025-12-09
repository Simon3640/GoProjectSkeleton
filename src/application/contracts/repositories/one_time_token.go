package contracts_repositories

import (
	dtos "github.com/simon3640/goprojectskeleton/src/application/shared/DTOs"
	application_errors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
)

type IOneTimeTokenRepository interface {
	IRepositoryBase[dtos.OneTimeTokenCreate, dtos.OneTimeTokenUpdate, models.OneTimeToken, models.OneTimeToken]
	GetByTokenHash(tokenHash []byte) (*models.OneTimeToken, *application_errors.ApplicationError)
}
