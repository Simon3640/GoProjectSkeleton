package repositories

type ModelConverter[ModelCreate any, Model any, DBModel any] interface {
	toGormCreate(model ModelCreate) *DBModel
	toDomain(ormModel *DBModel) *Model
}
