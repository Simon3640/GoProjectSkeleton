package repositories

import (
	"fmt"

	contracts_providers "gormgoskeleton/src/application/contracts/providers"
	contracts_repositories "gormgoskeleton/src/application/contracts/repositories"
	application_errors "gormgoskeleton/src/application/shared/errors"
	domain_utils "gormgoskeleton/src/domain/utils"

	"gorm.io/gorm"
)

type RepositoryBase[CreateModel any, UpdateModel any, Model any, DBModel any] struct {
	DB             *gorm.DB
	logger         contracts_providers.ILoggerProvider
	modelConverter ModelConverter[CreateModel, UpdateModel, Model, DBModel]
}

var _ contracts_repositories.IRepositoryBase[any, any, any, any] = (*RepositoryBase[any, any, any, any])(nil)

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

func (rb *RepositoryBase[CreateModel, UpdateModel, Model, DBModel]) Create(entity CreateModel) (*Model, *application_errors.ApplicationError) {
	// Convertir a modelo de GORM
	_entity := rb.modelConverter.ToGormCreate(entity)
	if err := rb.DB.Create(_entity).Error; err != nil {
		appErr := MapOrmError(err)
		rb.logger.Debug("Error creating entity", appErr.ToError())
		return nil, appErr
	}
	rb.logger.Debug("Entity created successfully", _entity)
	// Convertir de nuevo a modelo de dominio
	return rb.modelConverter.ToDomain(_entity), nil
}

func (rb *RepositoryBase[CreateModel, UpdateModel, Model, DBModel]) GetByID(id uint) (*Model, *application_errors.ApplicationError) {
	var entity DBModel
	if err := rb.DB.First(&entity, id).Error; err != nil {
		appErr := MapOrmError(err)
		rb.logger.Debug("Error retrieving entity", appErr.ToError())
		return nil, appErr
	}
	rb.logger.Debug("Entity retrieved successfully", entity)
	return rb.modelConverter.ToDomain(&entity), nil
}

func (rb *RepositoryBase[CreateModel, UpdateModel, Model, DBModel]) Update(id uint, entity UpdateModel) (*Model, *application_errors.ApplicationError) {
	updateData := rb.modelConverter.ToGormUpdate(entity)

	if err := rb.DB.Model(new(DBModel)).Where("id = ?", id).Updates(updateData).Error; err != nil {
		appErr := MapOrmError(err)
		rb.logger.Debug("Error updating entity", appErr.ToError())
		return nil, appErr
	}
	rb.logger.Debug("Entity updated successfully", updateData)
	updatedEntity, _ := rb.GetByID(id)
	return updatedEntity, nil
}

func (rb *RepositoryBase[CreateModel, UpdateModel, Model, DBModel]) SoftDelete(id uint) *application_errors.ApplicationError {
	if err := rb.DB.Delete(new(DBModel), id).Error; err != nil {
		appErr := MapOrmError(err)
		rb.logger.Debug("Error deleting entity", appErr.ToError())
		return appErr
	}
	return nil
}

func (rb *RepositoryBase[CreateModel, UpdateModel, Model, DBModel]) Delete(id uint) *application_errors.ApplicationError {
	if err := rb.DB.Unscoped().Delete(new(DBModel), id).Error; err != nil {
		appErr := MapOrmError(err)
		rb.logger.Debug("Error hard deleting entity", appErr.ToError())
		return appErr
	}
	return nil
}

func (rb *RepositoryBase[CreateModel, UpdateModel, Model, DBModel]) GetAll(payload *domain_utils.QueryPayloadBuilder[Model], skip int, limit int) ([]Model, *application_errors.ApplicationError) {
	var entities []DBModel
	// Apply filters from payload

	query := rb.DB.Model(new(DBModel))
	if payload != nil {
		for _, filter := range payload.Filters {
			gormCondition, args, err := FilterToGorm(filter)
			if err != nil {
				rb.logger.Debug("Error converting filter to GORM", err)
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
	query = query.Offset(skip).Limit(limit)
	// Execute the query
	if err := query.Find(&entities).Error; err != nil {
		if appErr := MapOrmError(err); appErr != nil {
			rb.logger.Debug("Error retrieving entities", appErr.ToError())
			return nil, appErr
		}
		return nil, DefaultORMError
	}

	result := make([]Model, len(entities))

	for i, entity := range entities {
		result[i] = *rb.modelConverter.ToDomain(&entity)
	}

	return result, nil
}
