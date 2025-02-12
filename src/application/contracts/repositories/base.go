package contracts_repositories

type IRepositoryBase[CreateDomainModel any, UpdateDomainModel any, DomainModel any, DBModel any] interface {
	Create(entity CreateDomainModel) (*DomainModel, error)
	GetByID(id int) (*DomainModel, error)
	Update(entity UpdateDomainModel) error
	Delete(id int) error
	GetAll() ([]DomainModel, error)
}
