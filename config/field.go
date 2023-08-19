package config

import (
	"reflect"

	"github.com/spf13/viper"
)

var Fields = make(map[string]entry)

type entry struct {
	Key         string
	Description string
	Default     any
	Marshal     func(any) (any, error)
	Unmarshal   func(any) (any, error)
	Validate    func(any) error
	SetValue    func(any) error
}

type registered[Raw, Value any] struct {
	Field[Raw, Value]
	value Value
}

func (r *registered[Raw, Value]) Get() Value {
	return r.value
}

func (r *registered[Raw, Value]) Set(value Value) error {
	marshalled, err := r.Marshal(value)
	if err != nil {
		return err
	}

	r.value = value
	viper.Set(r.Key, marshalled)
	return nil
}

func reg[Raw, Value any](field Field[Raw, Value]) *registered[Raw, Value] {
	if field.Marshal == nil {
		field.Marshal = func(value Value) (raw Raw, err error) {
			return reflect.
				ValueOf(value).
				Convert(reflect.ValueOf(raw).Type()).
				Interface().(Raw), nil
		}
	}
	if field.Unmarshal == nil {
		field.Unmarshal = func(raw Raw) (value Value, err error) {
			return reflect.
				ValueOf(raw).
				Convert(reflect.ValueOf(value).Type()).
				Interface().(Value), nil
		}
	}
	if field.Validate == nil {
		field.Validate = func(Value) error {
			return nil
		}
	}

	r := &registered[Raw, Value]{
		Field: field,
		value: field.Default,
	}

	Fields[field.Key] = entry{
		Key:         field.Key,
		Description: field.Description,
		Default:     field.Default,
		SetValue: func(a any) error {
			return r.Set(a.(Value))
		},
		Marshal: func(a any) (any, error) {
			return field.Marshal(a.(Value))
		},
		Unmarshal: func(a any) (any, error) {
			return field.Unmarshal(a.(Raw))
		},
		Validate: func(a any) error {
			return field.Validate(a.(Value))
		},
	}

	return r
}

type Field[Raw, Value any] struct {
	Key         string
	Description string
	Default     Value
	Validate    func(Value) error
	Unmarshal   func(Raw) (Value, error)
	Marshal     func(Value) (Raw, error)
}
