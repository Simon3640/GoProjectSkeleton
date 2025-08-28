package mocks

import (
	contracts_repositories "gormgoskeleton/src/application/contracts/repositories"
	dtos "gormgoskeleton/src/application/shared/DTOs"
	"gormgoskeleton/src/domain/models"
)

type MockOneTimeTokenRepository struct {
	MockRepositoryBase[dtos.OneTimeTokenCreate, dtos.OneTimeTokenUpdate, models.OneTimeToken, models.OneTimeToken]
}

var _ contracts_repositories.IOneTimeTokenRepository = (*MockOneTimeTokenRepository)(nil)
