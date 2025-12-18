// Package authrepositories contains the repository for the one time token model
package authrepositories

import (
	contractsproviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
	contractsrepositories "github.com/simon3640/goprojectskeleton/src/application/contracts/repositories"
	dtos "github.com/simon3640/goprojectskeleton/src/application/shared/DTOs"
	applicationerrors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
	dbmodels "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/models"
	reposhared "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/repositories/shared"

	"gorm.io/gorm"
)

// OneTimeTokenRepository is the repository for the one time token model
type OneTimeTokenRepository struct {
	reposhared.RepositoryBase[dtos.OneTimeTokenCreate, dtos.OneTimeTokenUpdate, models.OneTimeToken, dbmodels.OneTimeToken]
}

var _ contractsrepositories.IOneTimeTokenRepository = (*OneTimeTokenRepository)(nil)

// GetByTokenHash retrieves a one time token by its hash
func (or *OneTimeTokenRepository) GetByTokenHash(tokenHash []byte) (*models.OneTimeToken, *applicationerrors.ApplicationError) {
	var ormModel dbmodels.OneTimeToken

	if err := or.DB.Where("hash = ?", tokenHash).First(&ormModel).Error; err != nil {
		or.Logger.Debug("Error fetching one-time token by hash", err)
		return nil, reposhared.MapOrmError(err)
	}
	return or.ModelConverter.ToDomain(&ormModel), nil
}

// OneTimeTokenConverter is the converter for the one time token model
type OneTimeTokenConverter struct{}

var _ reposhared.ModelConverter[dtos.OneTimeTokenCreate, dtos.OneTimeTokenUpdate, models.OneTimeToken, dbmodels.OneTimeToken] = (*OneTimeTokenConverter)(nil)

// ToGormCreate converts a one time token create model to a one time token gorm model
func (uc *OneTimeTokenConverter) ToGormCreate(model dtos.OneTimeTokenCreate) *dbmodels.OneTimeToken {
	return &dbmodels.OneTimeToken{
		Purpose: string(model.Purpose),
		Hash:    model.Hash,
		UserID:  model.UserID,
		Expires: model.Expires,
		IsUsed:  false,
	}
}

// ToDomain converts a one time token gorm model to a one time token domain model
func (uc *OneTimeTokenConverter) ToDomain(ormModel *dbmodels.OneTimeToken) *models.OneTimeToken {
	return &models.OneTimeToken{
		DBBaseModel: models.DBBaseModel{
			ID:        ormModel.ID,
			CreatedAt: ormModel.CreatedAt,
			UpdatedAt: ormModel.UpdatedAt,
			DeletedAt: ormModel.DeletedAt.Time,
		},
		OneTimeTokenBase: models.OneTimeTokenBase{
			Purpose: models.OneTimeTokenPurpose(ormModel.Purpose),
			Hash:    ormModel.Hash,
			UserID:  ormModel.UserID,
			Expires: ormModel.Expires,
			IsUsed:  ormModel.IsUsed,
		},
	}
}

// ToGormUpdate converts a one time token update model to a one time token gorm model
func (uc *OneTimeTokenConverter) ToGormUpdate(model dtos.OneTimeTokenUpdate) *dbmodels.OneTimeToken {
	OneTimeToken := &dbmodels.OneTimeToken{}

	OneTimeToken.IsUsed = model.IsUsed
	OneTimeToken.ID = model.ID
	return OneTimeToken
}

// NewOneTimeTokenRepository creates a new one time token repository
func NewOneTimeTokenRepository(db *gorm.DB, logger contractsproviders.ILoggerProvider) *OneTimeTokenRepository {
	return &OneTimeTokenRepository{
		RepositoryBase: reposhared.RepositoryBase[
			dtos.OneTimeTokenCreate,
			dtos.OneTimeTokenUpdate,
			models.OneTimeToken,
			dbmodels.OneTimeToken,
		]{
			DB:             db,
			ModelConverter: &OneTimeTokenConverter{},
			Logger:         logger,
		},
	}
}
