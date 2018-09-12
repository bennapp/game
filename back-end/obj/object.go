package obj

import (
	"../typ"
	"github.com/google/uuid"
)

type Object struct {
	typ.Type
	Id uuid.UUID
}

type Objectable interface {
	ObjectId() string
}

func (object *Object) ObjectId() string {
	return object.Id.String()
}

func newObject(objectType string) Object {
	uuid, _ := uuid.NewUUID()
	return Object{Type: typ.NewType(objectType), Id: uuid}
}
