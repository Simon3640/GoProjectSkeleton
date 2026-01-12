package authcontracts

import (
	contractrepositories "github.com/simon3640/goprojectskeleton/src/application/contracts/repositories"
	authdtos "github.com/simon3640/goprojectskeleton/src/application/modules/auth/dtos"
	applicationerrors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	sharedmodels "github.com/simon3640/goprojectskeleton/src/domain/shared/models"
)

// IOneTimePasswordRepository is the interface for the one time password repository
type IOneTimePasswordRepository interface {
	contractrepositories.IRepositoryBase[authdtos.OneTimePasswordCreate, authdtos.OneTimePasswordUpdate, sharedmodels.OneTimePassword, sharedmodels.OneTimePassword]
	GetByPasswordHash(tokenHash []byte) (*sharedmodels.OneTimePassword, *applicationerrors.ApplicationError)
}
