package util

type Option[T any] struct {
	value  T
	isSome bool
}

func (o Option[T]) Unwrap() T {
	return o.value
}

func (o Option[T]) IsSome() bool {
	return o.isSome
}

func (o Option[T]) IsNone() bool {
	return !o.isSome
}

func Some[T any](value T) Option[T] {
	return Option[T]{value: value, isSome: true}
}

func None[T any]() Option[T] {
	return Option[T]{isSome: false}
}
