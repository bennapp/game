package rc

import (
	"../gs"
	"encoding/json"
	"fmt"
)

const OBJECT_LOCATION_LAYER = "o"
const PAINT_LOCATION_LAYER = "p"

type RedisStore interface {
	Key() string
	Value() string
}

type TypeStore struct {
	Type string
}

func newTypeStore(serializedData []byte) *TypeStore {
	typeStore := new(TypeStore)

	json.Unmarshal(serializedData, typeStore)

	return typeStore
}

type ObjectLocationStore struct {
	Layer    string
	Coord    gs.Coord
	ObjectId string
}

func newObjectLocationStoreRetriever(coord gs.Coord) *ObjectLocationStore {
	return &ObjectLocationStore{Layer: OBJECT_LOCATION_LAYER, Coord: coord}
}

func newObjectLocationStore(coord gs.Coord, object Dbo) *ObjectLocationStore {
	return &ObjectLocationStore{Layer: OBJECT_LOCATION_LAYER, Coord: coord, ObjectId: object.Id()}
}

func (olb *ObjectLocationStore) Key() string {
	return fmt.Sprintf("%v,%v", olb.Coord.X, olb.Coord.Y)
}

func (olb *ObjectLocationStore) Value() string {
	return fmt.Sprintf("%v:%v", olb.Layer, olb.ObjectId)
}

type ObjectStore struct {
	ObjectId         string
	SerializedObject []byte
}

func newObjectStoreRetriever(objectId string) *ObjectStore {
	return &ObjectStore{ObjectId: objectId}
}

func newObjectStore(object Dbo) *ObjectStore {
	serializedObject, _ := json.Marshal(object)

	return &ObjectStore{ObjectId: object.Id(), SerializedObject: serializedObject}
}

func (olb *ObjectStore) Key() string {
	return olb.ObjectId
}

func (olb *ObjectStore) Value() string {
	return string(olb.SerializedObject)
}

type PaintLocationStore struct {
	Layer string
	Coord gs.Coord
	Paint gs.Paint
}

func newPaintLocationStore(coord gs.Coord, paint gs.Paint) *PaintLocationStore {
	return &PaintLocationStore{Layer: PAINT_LOCATION_LAYER, Coord: coord, Paint: paint}
}
