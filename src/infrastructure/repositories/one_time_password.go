package repositories

import (
	contractsProviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
	contracts_repositories "github.com/simon3640/goprojectskeleton/src/application/contracts/repositories"
	dtos "github.com/simon3640/goprojectskeleton/src/application/shared/DTOs"
	application_errors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
	dbModels "github.com/simon3640/goprojectskeleton/src/infrastructure/database/goprojectskeleton/models"

	"gorm.io/gorm"
)

type OneTimePasswordRepository struct {
	RepositoryBase[dtos.OneTimePasswordCreate, dtos.OneTimePasswordUpdate, models.OneTimePassword, dbModels.OneTimePassword]
}

var _ contracts_repositories.IOneTimePasswordRepository = (*OneTimePasswordRepository)(nil)

func (or *OneTimePasswordRepository) GetByPasswordHash(tokenHash []byte) (*models.OneTimePassword, *application_errors.ApplicationError) {
	var ormModel dbModels.OneTimePassword

	if err := or.DB.Where("hash = ?", tokenHash).First(&ormModel).Error; err != nil {
		or.logger.Debug("Error fetching one-time token by hash", err)
		return nil, MapOrmError(err)
	}
	return or.modelConverter.ToDomain(&ormModel), nil
}

type OneTimePasswordConverter struct{}

var _ ModelConverter[dtos.OneTimePasswordCreate, dtos.OneTimePasswordUpdate, models.OneTimePassword, dbModels.OneTimePassword] = (*OneTimePasswordConverter)(nil)

func (uc *OneTimePasswordConverter) ToGormCreate(model dtos.OneTimePasswordCreate) *dbModels.OneTimePassword {
	return &dbModels.OneTimePassword{
		Purpose: string(model.Purpose),
		Hash:    model.Hash,
		UserID:  model.UserID,
		Expires: model.Expires,
		IsUsed:  false,
	}
}

func (uc *OneTimePasswordConverter) ToDomain(ormModel *dbModels.OneTimePassword) *models.OneTimePassword {
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

func (uc *OneTimePasswordConverter) ToGormUpdate(model dtos.OneTimePasswordUpdate) *dbModels.OneTimePassword {
	OneTimePassword := &dbModels.OneTimePassword{}

	OneTimePassword.IsUsed = model.IsUsed
	OneTimePassword.ID = model.ID
	return OneTimePassword
}

func NewOneTimePasswordRepository(db *gorm.DB, logger contractsProviders.ILoggerProvider) *OneTimePasswordRepository {
	return &OneTimePasswordRepository{
		RepositoryBase: RepositoryBase[
			dtos.OneTimePasswordCreate,
			dtos.OneTimePasswordUpdate,
			models.OneTimePassword,
			dbModels.OneTimePassword,
		]{DB: db, modelConverter: &OneTimePasswordConverter{}, logger: logger},
	}
}
