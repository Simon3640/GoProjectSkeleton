package repositories

import (
	contracts_providers "gormgoskeleton/src/application/contracts/providers"
	contracts_repositories "gormgoskeleton/src/application/contracts/repositories"
	dtos "gormgoskeleton/src/application/shared/DTOs"
	application_errors "gormgoskeleton/src/application/shared/errors"
	"gormgoskeleton/src/domain/models"
	db_models "gormgoskeleton/src/infrastructure/database/gormgoskeleton/models"

	"gorm.io/gorm"
)

type OneTimePasswordRepository struct {
	RepositoryBase[dtos.OneTimePasswordCreate, dtos.OneTimePasswordUpdate, models.OneTimePassword, db_models.OneTimePassword]
}

var _ contracts_repositories.IOneTimePasswordRepository = (*OneTimePasswordRepository)(nil)

func (or *OneTimePasswordRepository) GetByPasswordHash(tokenHash []byte) (*models.OneTimePassword, *application_errors.ApplicationError) {
	var ormModel db_models.OneTimePassword

	if err := or.DB.Where("hash = ?", tokenHash).First(&ormModel).Error; err != nil {
		or.logger.Debug("Error fetching one-time token by hash", err)
		return nil, MapOrmError(err)
	}
	return or.modelConverter.ToDomain(&ormModel), nil
}

type OneTimePasswordConverter struct{}

var _ ModelConverter[dtos.OneTimePasswordCreate, dtos.OneTimePasswordUpdate, models.OneTimePassword, db_models.OneTimePassword] = (*OneTimePasswordConverter)(nil)

func (uc *OneTimePasswordConverter) ToGormCreate(model dtos.OneTimePasswordCreate) *db_models.OneTimePassword {
	return &db_models.OneTimePassword{
		Purpose: string(model.Purpose),
		Hash:    model.Hash,
		UserID:  model.UserID,
		Expires: model.Expires,
		IsUsed:  false,
	}
}

func (uc *OneTimePasswordConverter) ToDomain(ormModel *db_models.OneTimePassword) *models.OneTimePassword {
	return &models.OneTimePassword{
		DBBaseModel: models.DBBaseModel{
			ID:        ormModel.ID,
			CreatedAt: ormModel.CreatedAt,
			UpdatedAt: ormModel.UpdatedAt,
			DeletedAt: ormModel.DeletedAt.Time,
		},
		OneTimePasswordBase: models.OneTimePasswordBase{
			Purpose: models.OneTimePasswordPurpose(ormModel.Purpose),
			Hash:    ormModel.Hash,
			UserID:  ormModel.UserID,
			Expires: ormModel.Expires,
			IsUsed:  ormModel.IsUsed,
		},
	}
}

func (uc *OneTimePasswordConverter) ToGormUpdate(model dtos.OneTimePasswordUpdate) *db_models.OneTimePassword {
	OneTimePassword := &db_models.OneTimePassword{}

	OneTimePassword.IsUsed = model.IsUsed
	OneTimePassword.ID = model.ID
	return OneTimePassword
}

func NewOneTimePasswordRepository(db *gorm.DB, logger contracts_providers.ILoggerProvider) *OneTimePasswordRepository {
	return &OneTimePasswordRepository{
		RepositoryBase: RepositoryBase[
			dtos.OneTimePasswordCreate,
			dtos.OneTimePasswordUpdate,
			models.OneTimePassword,
			db_models.OneTimePassword,
		]{DB: db, modelConverter: &OneTimePasswordConverter{}, logger: logger},
	}
}
