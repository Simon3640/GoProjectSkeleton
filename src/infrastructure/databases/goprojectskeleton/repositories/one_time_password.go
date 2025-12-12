package repositories

import (
	contractsProviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
	authcontracts "github.com/simon3640/goprojectskeleton/src/application/modules/auth/contracts"
	authdtos "github.com/simon3640/goprojectskeleton/src/application/modules/auth/dtos"
	application_errors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
	dbmodels "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/models"

	"gorm.io/gorm"
)

// OneTimePasswordRepository is the repository for the one time password model
type OneTimePasswordRepository struct {
	RepositoryBase[authdtos.OneTimePasswordCreate, authdtos.OneTimePasswordUpdate, models.OneTimePassword, dbmodels.OneTimePassword]
}

var _ authcontracts.IOneTimePasswordRepository = (*OneTimePasswordRepository)(nil)

// GetByPasswordHash retrieves a one time password by its hash
func (or *OneTimePasswordRepository) GetByPasswordHash(tokenHash []byte) (*models.OneTimePassword, *application_errors.ApplicationError) {
	var ormModel dbmodels.OneTimePassword

	if err := or.DB.Where("hash = ?", tokenHash).First(&ormModel).Error; err != nil {
		or.logger.Debug("Error fetching one-time token by hash", err)
		return nil, MapOrmError(err)
	}
	return or.modelConverter.ToDomain(&ormModel), nil
}

// OneTimePasswordConverter is the converter for the one time password model
type OneTimePasswordConverter struct{}

var _ ModelConverter[authdtos.OneTimePasswordCreate, authdtos.OneTimePasswordUpdate, models.OneTimePassword, dbmodels.OneTimePassword] = (*OneTimePasswordConverter)(nil)

// ToGormCreate converts a one time password create model to a one time password gorm model
func (uc *OneTimePasswordConverter) ToGormCreate(model authdtos.OneTimePasswordCreate) *dbmodels.OneTimePassword {
	return &dbmodels.OneTimePassword{
		Purpose: string(model.Purpose),
		Hash:    model.Hash,
		UserID:  model.UserID,
		Expires: model.Expires,
		IsUsed:  false,
	}
}

// ToDomain converts a one time password gorm model to a one time password domain model
func (uc *OneTimePasswordConverter) ToDomain(ormModel *dbmodels.OneTimePassword) *models.OneTimePassword {
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

// ToGormUpdate converts a one time password update model to a one time password gorm model
func (uc *OneTimePasswordConverter) ToGormUpdate(model authdtos.OneTimePasswordUpdate) *dbmodels.OneTimePassword {
	OneTimePassword := &dbmodels.OneTimePassword{}

	OneTimePassword.IsUsed = model.IsUsed
	OneTimePassword.ID = model.ID
	return OneTimePassword
}

// NewOneTimePasswordRepository creates a new one time password repository
func NewOneTimePasswordRepository(db *gorm.DB, logger contractsProviders.ILoggerProvider) *OneTimePasswordRepository {
	return &OneTimePasswordRepository{
		RepositoryBase: RepositoryBase[
			authdtos.OneTimePasswordCreate,
			authdtos.OneTimePasswordUpdate,
			models.OneTimePassword,
			dbmodels.OneTimePassword,
		]{DB: db, modelConverter: &OneTimePasswordConverter{}, logger: logger},
	}
}
