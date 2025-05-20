package db_models

type DBModel interface {
	TableName() string
}
