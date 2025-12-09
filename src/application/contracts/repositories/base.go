package contracts_repositories

import (
	application_errors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	domain_utils "github.com/simon3640/goprojectskeleton/src/domain/utils"
)

type IRepositoryBase[CreateDomainModel any, UpdateDomainModel any, DomainModel any, DBModel any] interface {
	Create(entity CreateDomainModel) (*DomainModel, *application_errors.ApplicationError)
	GetByID(id uint) (*DomainModel, *application_errors.ApplicationError)
	Update(id uint, entity UpdateDomainModel) (*DomainModel, *application_errors.ApplicationError)
	Delete(id uint) *application_errors.ApplicationError
	SoftDelete(id uint) *application_errors.ApplicationError
	GetAll(payload *domain_utils.QueryPayloadBuilder[DomainModel], skip int, limit int) ([]DomainModel, int64, *application_errors.ApplicationError)
}
