package el

import "github.com/google/uuid"

type Object struct {
	Id   uuid.UUID
	Type string
}

func (object *Object) DboId() string {
	return object.Id.String()
}

func newObject(objectType string) Object {
	uuid, _ := uuid.NewUUID()
	return Object{Type: objectType, Id: uuid}
}

func loadObject(objectType string) Object {
	return Object{Type: objectType}
}
