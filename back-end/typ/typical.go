package typ

type Typical interface {
	GetType() string
}

type Type struct {
	Type string
}

func (typical *Type) GetType() string {
	return typical.Type
}

func NewType(typeString string) Type {
	return Type{Type: typeString}
}
