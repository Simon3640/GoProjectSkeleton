package contracts_repositories

import (
	dtos "github.com/simon3640/goprojectskeleton/src/application/shared/DTOs"
	application_errors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	sharedmodels "github.com/simon3640/goprojectskeleton/src/domain/shared/models"
)

type IOneTimeTokenRepository interface {
	IRepositoryBase[dtos.OneTimeTokenCreate, dtos.OneTimeTokenUpdate, sharedmodels.OneTimeToken, sharedmodels.OneTimeToken]
	GetByTokenHash(tokenHash []byte) (*sharedmodels.OneTimeToken, *application_errors.ApplicationError)
}
