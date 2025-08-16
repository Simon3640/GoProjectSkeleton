package repositories

import (
	"gormgoskeleton/src/application/contracts"
	contracts_repositories "gormgoskeleton/src/application/contracts/repositories"
	"gormgoskeleton/src/domain/models"
	db_models "gormgoskeleton/src/infrastructure/database/gormgoskeleton/models"

	"gorm.io/gorm"
)

type RoleRepository struct {
	RepositoryBase[models.RoleCreate, models.RoleUpdate, models.Role, db_models.Role]
}

var _ contracts_repositories.IRoleRepository = (*RoleRepository)(nil)

type RoleConverter struct{}

var _ ModelConverter[models.RoleCreate, models.RoleUpdate, models.Role, db_models.Role] = (*RoleConverter)(nil)

func (uc *RoleConverter) ToGormCreate(model models.RoleCreate) *db_models.Role {
	return &db_models.Role{
		Key:      model.Key,
		IsActive: model.IsActive,
		Priority: model.Priority,
	}
}

func (uc *RoleConverter) ToDomain(ormModel *db_models.Role) *models.Role {
	return &models.Role{
		ID: ormModel.ID,
		RoleBase: models.RoleBase{
			Key:      ormModel.Key,
			IsActive: ormModel.IsActive,
			Priority: ormModel.Priority,
		},
	}
}

func (uc *RoleConverter) ToGormUpdate(model models.RoleUpdate) *db_models.Role {
	Role := &db_models.Role{}

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

func NewRoleRepository(db *gorm.DB, logger contracts.ILoggerProvider) *RoleRepository {
	return &RoleRepository{
		RepositoryBase: RepositoryBase[
			models.RoleCreate,
			models.RoleUpdate,
			models.Role,
			db_models.Role,
		]{DB: db, modelConverter: &RoleConverter{}, logger: logger},
	}
}
