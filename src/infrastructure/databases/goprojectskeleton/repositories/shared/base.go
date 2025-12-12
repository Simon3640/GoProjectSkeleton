// Package shared contains the base repository for the database models
package shared

import (
	"fmt"

	contractsproviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
	contractsrepositories "github.com/simon3640/goprojectskeleton/src/application/contracts/repositories"
	applicationerrors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	domain_utils "github.com/simon3640/goprojectskeleton/src/domain/utils"

	"gorm.io/gorm"
)

// RepositoryBase is the abstract base repository for the database models
// It contains the database connection, the logger and the model converter
// It implements the IRepositoryBase interface
type RepositoryBase[CreateModel any, UpdateModel any, Model any, DBModel any] struct {
	DB             *gorm.DB
	Logger         contractsproviders.ILoggerProvider
	ModelConverter ModelConverter[CreateModel, UpdateModel, Model, DBModel]
}

func SetUpRepositoryBase[CreateModel, UpdateModel, Model, DBModel any](db *gorm.DB,
	logger contractsproviders.ILoggerProvider,
	modelConverter ModelConverter[CreateModel, UpdateModel, Model, DBModel],
) RepositoryBase[CreateModel, UpdateModel, Model, DBModel] {
	return RepositoryBase[CreateModel, UpdateModel, Model, DBModel]{
		DB:             db,
		Logger:         logger,
		ModelConverter: modelConverter,
	}
}

var _ contractsrepositories.IRepositoryBase[any, any, any, any] = (*RepositoryBase[any, any, any, any])(nil)

// FilterToGorm converts a filter to a GORM filter
func FilterToGorm(f domain_utils.Filter) (string, []interface{}, error) {
	switch f.Operator {
	case domain_utils.OperatorEqual:
		return f.Field + " = ?", []interface{}{*f.Value}, nil
	case domain_utils.OperatorNotEqual:
		return f.Field + " != ?", []interface{}{*f.Value}, nil
	case domain_utils.OperatorGreaterThan:
		return f.Field + " > ?", []interface{}{*f.Value}, nil
	case domain_utils.OperatorGreaterEqual:
		return f.Field + " >= ?", []interface{}{*f.Value}, nil
	case domain_utils.OperatorLessThan:
		return f.Field + " < ?", []interface{}{*f.Value}, nil
	case domain_utils.OperatorLessEqual:
		return f.Field + " <= ?", []interface{}{*f.Value}, nil
	case domain_utils.OperatorLike:
		return f.Field + " LIKE ?", []interface{}{*f.Value}, nil
	case domain_utils.OperatorILike:
		return f.Field + " ILIKE ?", []interface{}{*f.Value}, nil
	case domain_utils.OperatorIn:
		return f.Field + " IN ?", []interface{}{*f.Value}, nil
	case domain_utils.OperatorNotIn:
		return f.Field + " NOT IN ?", []interface{}{*f.Value}, nil
	case domain_utils.OperatorIsNull:
		return f.Field + " IS NULL", []interface{}{}, nil
	case domain_utils.OperatorIsNotNull:
		return f.Field + " IS NOT NULL", []interface{}{}, nil
	default:
		return "", nil, fmt.Errorf("unsupported operator: %s", f.Operator)
	}
}

// SortToGorm converts a sort to a GORM sort
func SortToGorm(s domain_utils.Sort) string {
	switch s.Field {
	case "CreatedAt":
		return fmt.Sprintf("created_at %s", s.Order)
	case "UpdatedAt":
		return fmt.Sprintf("updated_at %s", s.Order)
	case "DeletedAt":
		return fmt.Sprintf("deleted_at %s", s.Order)
	default:
		return fmt.Sprintf("%s %s", s.Field, s.Order)
	}
}

// Create creates a new entity
func (rb *RepositoryBase[CreateModel, UpdateModel, Model, DBModel]) Create(entity CreateModel) (*Model, *applicationerrors.ApplicationError) {
	// Convertir a modelo de GORM
	_entity := rb.ModelConverter.ToGormCreate(entity)
	if err := rb.DB.Create(_entity).Error; err != nil {
		appErr := MapOrmError(err)
		rb.Logger.Debug("Error creating entity", appErr.ToError())
		return nil, appErr
	}
	rb.Logger.Debug("Entity created successfully", _entity)
	// Convertir de nuevo a modelo de dominio
	return rb.ModelConverter.ToDomain(_entity), nil
}

// GetByID retrieves an entity by its ID
func (rb *RepositoryBase[CreateModel, UpdateModel, Model, DBModel]) GetByID(id uint) (*Model, *applicationerrors.ApplicationError) {
	var entity DBModel
	if err := rb.DB.First(&entity, id).Error; err != nil {
		appErr := MapOrmError(err)
		rb.Logger.Debug("Error retrieving entity", appErr.ToError())
		return nil, appErr
	}
	rb.Logger.Debug("Entity retrieved successfully", entity)
	return rb.ModelConverter.ToDomain(&entity), nil
}

// Update updates an entity
func (rb *RepositoryBase[CreateModel, UpdateModel, Model, DBModel]) Update(id uint, entity UpdateModel) (*Model, *applicationerrors.ApplicationError) {
	updateData := rb.ModelConverter.ToGormUpdate(entity)

	if err := rb.DB.Model(new(DBModel)).Where("id = ?", id).Updates(updateData).Error; err != nil {
		appErr := MapOrmError(err)
		rb.Logger.Debug("Error updating entity", appErr.ToError())
		return nil, appErr
	}
	rb.Logger.Debug("Entity updated successfully", updateData)
	updatedEntity, _ := rb.GetByID(id)
	return updatedEntity, nil
}

// SoftDelete soft deletes an entity
func (rb *RepositoryBase[CreateModel, UpdateModel, Model, DBModel]) SoftDelete(id uint) *applicationerrors.ApplicationError {
	if err := rb.DB.Delete(new(DBModel), id).Error; err != nil {
		appErr := MapOrmError(err)
		rb.Logger.Debug("Error deleting entity", appErr.ToError())
		return appErr
	}
	return nil
}

// Delete hard deletes an entity
func (rb *RepositoryBase[CreateModel, UpdateModel, Model, DBModel]) Delete(id uint) *applicationerrors.ApplicationError {
	if err := rb.DB.Unscoped().Delete(new(DBModel), id).Error; err != nil {
		appErr := MapOrmError(err)
		rb.Logger.Debug("Error hard deleting entity", appErr.ToError())
		return appErr
	}
	return nil
}

// GetAll retrieves all entities
func (rb *RepositoryBase[CreateModel, UpdateModel, Model, DBModel]) GetAll(payload *domain_utils.QueryPayloadBuilder[Model], skip, limit int) ([]Model, int64, *applicationerrors.ApplicationError) {
	var entities []DBModel
	// Apply filters from payload

	query := rb.DB.Model(new(DBModel))
	if payload != nil {
		for _, filter := range payload.Filters {
			gormCondition, args, err := FilterToGorm(filter)
			if err != nil {
				rb.Logger.Debug("Error converting filter to GORM", err)
				continue
			}
			if args != nil {
				query = query.Where(gormCondition, args...)
			} else {
				query = query.Where(gormCondition)
			}
		}
		// Apply sorts from payload
		for _, sort := range payload.Sorts {
			gormSort := SortToGorm(sort)
			query = query.Order(gormSort)
		}
	}

	// Count total records
	var total int64
	if err := query.Count(&total).Error; err != nil {
		if appErr := MapOrmError(err); appErr != nil {
			rb.Logger.Debug("Error counting entities", appErr.ToError())
			return nil, 0, appErr
		}
		return nil, 0, DefaultORMError
	}

	query = query.Offset(skip).Limit(limit)
	// Execute the query
	if err := query.Find(&entities).Error; err != nil {
		if appErr := MapOrmError(err); appErr != nil {
			rb.Logger.Debug("Error retrieving entities", appErr.ToError())
			return nil, 0, appErr
		}
		return nil, 0, DefaultORMError
	}

	result := make([]Model, len(entities))

	for i, entity := range entities {
		result[i] = *rb.ModelConverter.ToDomain(&entity)
	}

	return result, total, nil
}
