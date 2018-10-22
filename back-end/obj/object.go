package obj

import (
	"../typ"
	"fmt"
	"github.com/google/uuid"
)

type Object struct {
	typ.Type
	UUID uuid.UUID
	Id   string
}

type Objectable interface {
	ObjectId() string
}

func (object *Object) ObjectId() string {
	return object.UUID.String()
}

func newObject(objectType string) Object {
	uuid, _ := uuid.NewUUID()
	return Object{Type: typ.NewType(objectType), UUID: uuid, Id: fmt.Sprintf("%v", uuid)}
}
