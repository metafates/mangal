// Code generated by "enumer -type=Type -trimprefix=Type -json -text"; DO NOT EDIT.

package info

import (
	"encoding/json"
	"fmt"
	"strings"
)

const _TypeName = "BundleLua"

var _TypeIndex = [...]uint8{0, 6, 9}

const _TypeLowerName = "bundlelua"

func (i Type) String() string {
	i -= 1
	if i >= Type(len(_TypeIndex)-1) {
		return fmt.Sprintf("Type(%d)", i+1)
	}
	return _TypeName[_TypeIndex[i]:_TypeIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _TypeNoOp() {
	var x [1]struct{}
	_ = x[TypeBundle-(1)]
	_ = x[TypeLua-(2)]
}

var _TypeValues = []Type{TypeBundle, TypeLua}

var _TypeNameToValueMap = map[string]Type{
	_TypeName[0:6]:      TypeBundle,
	_TypeLowerName[0:6]: TypeBundle,
	_TypeName[6:9]:      TypeLua,
	_TypeLowerName[6:9]: TypeLua,
}

var _TypeNames = []string{
	_TypeName[0:6],
	_TypeName[6:9],
}

// TypeString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func TypeString(s string) (Type, error) {
	if val, ok := _TypeNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _TypeNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to Type values", s)
}

// TypeValues returns all values of the enum
func TypeValues() []Type {
	return _TypeValues
}

// TypeStrings returns a slice of all String values of the enum
func TypeStrings() []string {
	strs := make([]string, len(_TypeNames))
	copy(strs, _TypeNames)
	return strs
}

// IsAType returns "true" if the value is listed in the enum definition. "false" otherwise
func (i Type) IsAType() bool {
	for _, v := range _TypeValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for Type
func (i Type) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for Type
func (i *Type) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("Type should be a string, got %s", data)
	}

	var err error
	*i, err = TypeString(s)
	return err
}

// MarshalText implements the encoding.TextMarshaler interface for Type
func (i Type) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for Type
func (i *Type) UnmarshalText(text []byte) error {
	var err error
	*i, err = TypeString(string(text))
	return err
}