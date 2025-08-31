package contracts_repositories

import (
	dtos "gormgoskeleton/src/application/shared/DTOs"
	application_errors "gormgoskeleton/src/application/shared/errors"
	"gormgoskeleton/src/domain/models"
)

type IOneTimePasswordRepository interface {
	IRepositoryBase[dtos.OneTimePasswordCreate, dtos.OneTimePasswordUpdate, models.OneTimePassword, models.OneTimePassword]
	GetByPasswordHash(tokenHash []byte) (*models.OneTimePassword, *application_errors.ApplicationError)
}
