package dbpg

type ITable[T any] interface {
	New() T
	Clone() T
	TableName() string
}