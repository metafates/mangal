package icon

//go:generate enumer -type=Type -trimprefix=Type -json -yaml -text
type Type uint8

const (
	TypeNerd Type = iota + 1
	TypeASCII
)

var currentType Type = TypeASCII

func SetType(iconType Type) {
	currentType = iconType
}

func GetType() Type {
	return currentType
}
