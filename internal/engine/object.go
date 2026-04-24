package engine

const (
	LIST   = "list"
	STRING = "string"
)

type Object struct {
	typ   string
	value any
}

func NewObject(typ string, value any) *Object {
	return &Object{
		typ:   typ,
		value: value,
	}
}
