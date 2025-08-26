package contracts_repositories

import (
	application_errors "gormgoskeleton/src/application/shared/errors"
)

type IRepositoryBase[CreateDomainModel any, UpdateDomainModel any, DomainModel any, DBModel any] interface {
	Create(entity CreateDomainModel) (*DomainModel, *application_errors.ApplicationError)
	GetByID(id uint) (*DomainModel, *application_errors.ApplicationError)
	Update(id uint, entity UpdateDomainModel) (*DomainModel, *application_errors.ApplicationError)
	Delete(id uint) *application_errors.ApplicationError
	SoftDelete(id uint) *application_errors.ApplicationError
	GetAll(payload *map[string]string, skip *int, limit *int) ([]DomainModel, *application_errors.ApplicationError)
}
