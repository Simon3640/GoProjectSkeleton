package contracts_repositories

import (
	dtos "goprojectskeleton/src/application/shared/DTOs"
	application_errors "goprojectskeleton/src/application/shared/errors"
	"goprojectskeleton/src/domain/models"
)

type IOneTimeTokenRepository interface {
	IRepositoryBase[dtos.OneTimeTokenCreate, dtos.OneTimeTokenUpdate, models.OneTimeToken, models.OneTimeToken]
	GetByTokenHash(tokenHash []byte) (*models.OneTimeToken, *application_errors.ApplicationError)
}
