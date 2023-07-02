package icon

type Type uint8

const (
	typeNerd Type = iota + 1
	typeASCII
)

func SetType(iconType Type) {
}
