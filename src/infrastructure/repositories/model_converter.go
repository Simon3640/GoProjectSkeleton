package repositories

type ModelConverter[ModelCreate any, ModelUpdate any, Model any, DBModel any] interface {
	toGormCreate(model ModelCreate) *DBModel
	toDomain(ormModel *DBModel) *Model
	toGormUpdate(model ModelUpdate) *DBModel
}
