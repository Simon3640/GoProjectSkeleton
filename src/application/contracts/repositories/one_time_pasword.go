package contracts_repositories

import (
	dtos "github.com/simon3640/goprojectskeleton/src/application/shared/DTOs"
	application_errors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
)

type IOneTimePasswordRepository interface {
	IRepositoryBase[dtos.OneTimePasswordCreate, dtos.OneTimePasswordUpdate, models.OneTimePassword, models.OneTimePassword]
	GetByPasswordHash(tokenHash []byte) (*models.OneTimePassword, *application_errors.ApplicationError)
}
