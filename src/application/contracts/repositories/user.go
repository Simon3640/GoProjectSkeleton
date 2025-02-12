package contracts_repositories

// import "gormgoskeleton/src/domain/models"

type IUserRepository[ModelCreate any, ModelUpdate any, Model any, DBModel any] interface {
	IRepositoryBase[ModelCreate, ModelUpdate, Model, DBModel]
}
