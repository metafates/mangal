package util

// Option encapsulates an optional value.
// Similar to the Maybe monad in Haskell or Option in Rust.
type Option[T any] struct {
	value  T
	isSome bool
}

// Unwrap returns an internal value of the option.
// Panics if the option is None.
func (o Option[T]) Unwrap() T {
	if !o.IsSome() {
		panic("called `Option.Unwrap()` on a `None` value")
	}

	return o.value
}

// UnwrapOr returns an internal value of the option.
// Returns the default value if the option is None.
func (o Option[T]) UnwrapOr(or T) T {
	if !o.IsSome() {
		return or
	}

	return o.value
}

// IsSome returns true if the option is Some.
func (o Option[T]) IsSome() bool {
	return o.isSome
}

// IsNone returns true if the option is None.
func (o Option[T]) IsNone() bool {
	return !o.isSome
}

// Some returns a new Some option.
func Some[T any](value T) Option[T] {
	return Option[T]{value: value, isSome: true}
}

// None returns a new None option.
func None[T any]() Option[T] {
	return Option[T]{isSome: false}
}
