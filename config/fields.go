package config

type Field interface {
	Value() any
	Default() any

	Description() string
	Key() string
}

type field[T any] struct {
	key          string
	defaultValue T
	description  string
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

var Fields []Field

func newField[T any](key string, defaultValue T, description string) field[T] {
	f := field[T]{
		key:          key,
		defaultValue: defaultValue,
		description:  description,
	}

	Fields = append(Fields, f)

	return f
}

var (
	DownloadFormat = newField(
		"download.format",
		"pdf",
		"Format to download chapters in",
	)

	DownloadPath = newField(
		"download.path",
		".",
		"Path where chapters will be downloaded",
	)

	ReadFormat = newField(
		"read.format",
		"pdf",
		"Format to read chapters in",
	)
)
