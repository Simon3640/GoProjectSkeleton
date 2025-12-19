package userrepositories

import (
	contractsprovider "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
	usercontracts "github.com/simon3640/goprojectskeleton/src/application/modules/user/contracts"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
	dbmodels "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/models"
	reposhared "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/repositories/shared"

	"gorm.io/gorm"
)

// RoleRepository is the repository for the role model
type RoleRepository struct {
	reposhared.RepositoryBase[models.RoleCreate, models.RoleUpdate, models.Role, dbmodels.Role]
}

var _ usercontracts.IRoleRepository = (*RoleRepository)(nil)

// RoleConverter is the converter for the role model
type RoleConverter struct{}

var _ reposhared.ModelConverter[models.RoleCreate, models.RoleUpdate, models.Role, dbmodels.Role] = (*RoleConverter)(nil)

// ToGormCreate converts a role create model to a role gorm model
func (uc *RoleConverter) ToGormCreate(model models.RoleCreate) *dbmodels.Role {
	return &dbmodels.Role{
		Key:      model.Key,
		IsActive: model.IsActive,
		Priority: model.Priority,
	}
}

// ToDomain converts a role gorm model to a role domain model
func (uc *RoleConverter) ToDomain(ormModel *dbmodels.Role) *models.Role {
	return &models.Role{
		ID: ormModel.ID,
		RoleBase: models.RoleBase{
			Key:      ormModel.Key,
			IsActive: ormModel.IsActive,
			Priority: ormModel.Priority,
		},
	}
}

// ToGormUpdate converts a role update model to a role gorm model
func (uc *RoleConverter) ToGormUpdate(model models.RoleUpdate) *dbmodels.Role {
	Role := &dbmodels.Role{}

	if model.Key != nil {
		Role.Key = *model.Key
	}

	if model.IsActive != nil {
		Role.IsActive = *model.IsActive
	}

	if model.Priority != nil {
		Role.Priority = *model.Priority
	}

	Role.ID = model.ID
	return Role
}

// NewRoleRepository creates a new role repository
func NewRoleRepository(db *gorm.DB, logger contractsprovider.ILoggerProvider) *RoleRepository {
	return &RoleRepository{
		RepositoryBase: reposhared.RepositoryBase[
			models.RoleCreate,
			models.RoleUpdate,
			models.Role,
			dbmodels.Role,
		]{
			DB:             db,
			ModelConverter: &RoleConverter{},
			Logger:         logger,
		},
	}
}
