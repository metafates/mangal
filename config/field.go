package config

var Fields []Field

func register[T any](f field[T]) field[T] {
	Fields = append(Fields, &f)
	return f
}

type Field interface {
	Key() string
	Value() any
	Default() any
	Description() string

	Transform() (any, error)
	Init() error
}

type field[T any] struct {
	key          string
	defaultValue T
	description  string
	transform    func(T) (T, error)
	init         func(T) error
}

func (f field[T]) Default() any {
	return f.defaultValue
}

func (f field[T]) Description() string {
	return f.description
}

func (f field[T]) Key() string {
	return f.key
}

func (f field[T]) Value() any {
	return instance.Get(f.key)
}

func (f field[T]) Get() T {
	return f.Value().(T)
}

func (f field[T]) Init() error {
	if f.init == nil {
		return nil
	}

	return f.init(f.Get())
}

func (f field[T]) Transform() (any, error) {
	if f.transform == nil {
		return f.Get(), nil
	}

	return f.transform(f.Get())
}
