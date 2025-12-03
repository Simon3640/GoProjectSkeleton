package contracts_repositories

import (
	dtos "goprojectskeleton/src/application/shared/DTOs"
	application_errors "goprojectskeleton/src/application/shared/errors"
	"goprojectskeleton/src/domain/models"
)

type IOneTimePasswordRepository interface {
	IRepositoryBase[dtos.OneTimePasswordCreate, dtos.OneTimePasswordUpdate, models.OneTimePassword, models.OneTimePassword]
	GetByPasswordHash(tokenHash []byte) (*models.OneTimePassword, *application_errors.ApplicationError)
}
