package typ

type Type struct {
	Type string
}

type Typical interface {
	GetType() string
}

func (typical *Type) GetType() string {
	return typical.Type
}

func NewType(typeString string) Type {
	return Type{Type: typeString}
}
