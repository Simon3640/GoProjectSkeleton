package repositories

import (
	contractsProviders "gormgoskeleton/src/application/contracts/providers"
	contracts_repositories "gormgoskeleton/src/application/contracts/repositories"
	dtos "gormgoskeleton/src/application/shared/DTOs"
	application_errors "gormgoskeleton/src/application/shared/errors"
	"gormgoskeleton/src/domain/models"
	dbModels "gormgoskeleton/src/infrastructure/database/gormgoskeleton/models"

	"gorm.io/gorm"
)

type OneTimeTokenRepository struct {
	RepositoryBase[dtos.OneTimeTokenCreate, dtos.OneTimeTokenUpdate, models.OneTimeToken, dbModels.OneTimeToken]
}

var _ contracts_repositories.IOneTimeTokenRepository = (*OneTimeTokenRepository)(nil)

func (or *OneTimeTokenRepository) GetByTokenHash(tokenHash []byte) (*models.OneTimeToken, *application_errors.ApplicationError) {
	var ormModel dbModels.OneTimeToken

	if err := or.DB.Where("hash = ?", tokenHash).First(&ormModel).Error; err != nil {
		or.logger.Debug("Error fetching one-time token by hash", err)
		return nil, MapOrmError(err)
	}
	return or.modelConverter.ToDomain(&ormModel), nil
}

type OneTimeTokenConverter struct{}

var _ ModelConverter[dtos.OneTimeTokenCreate, dtos.OneTimeTokenUpdate, models.OneTimeToken, dbModels.OneTimeToken] = (*OneTimeTokenConverter)(nil)

func (uc *OneTimeTokenConverter) ToGormCreate(model dtos.OneTimeTokenCreate) *dbModels.OneTimeToken {
	return &dbModels.OneTimeToken{
		Purpose: string(model.Purpose),
		Hash:    model.Hash,
		UserID:  model.UserID,
		Expires: model.Expires,
		IsUsed:  false,
	}
}

func (uc *OneTimeTokenConverter) ToDomain(ormModel *dbModels.OneTimeToken) *models.OneTimeToken {
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

func (uc *OneTimeTokenConverter) ToGormUpdate(model dtos.OneTimeTokenUpdate) *dbModels.OneTimeToken {
	OneTimeToken := &dbModels.OneTimeToken{}

	OneTimeToken.IsUsed = model.IsUsed
	OneTimeToken.ID = model.ID
	return OneTimeToken
}

func NewOneTimeTokenRepository(db *gorm.DB, logger contractsProviders.ILoggerProvider) *OneTimeTokenRepository {
	return &OneTimeTokenRepository{
		RepositoryBase: RepositoryBase[
			dtos.OneTimeTokenCreate,
			dtos.OneTimeTokenUpdate,
			models.OneTimeToken,
			dbModels.OneTimeToken,
		]{DB: db, modelConverter: &OneTimeTokenConverter{}, logger: logger},
	}
}
