package domain

type IFactory[T any] interface {
	Create() (T, error)
}
