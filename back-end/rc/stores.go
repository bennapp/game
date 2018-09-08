package rc

import (
	"../gs"
	"../obj"
	"../pnt"
	"encoding/json"
	"fmt"
)

const OBJECT_LOCATION_LAYER = "o"
const PAINT_LOCATION_LAYER = "p"

type RedisStore interface {
	Key() string
	Value() string
}

type TypeDeserializer struct {
	Type string
}

func newTypeDeserializer(serializedData []byte) *TypeDeserializer {
	typeDeserializer := new(TypeDeserializer)

	json.Unmarshal(serializedData, typeDeserializer)

	return typeDeserializer
}

type ObjectLocationStore struct {
	Coord    gs.Coord
	ObjectId string
}

func newObjectLocationStoreRetriever(coord gs.Coord) *ObjectLocationStore {
	return &ObjectLocationStore{Coord: coord}
}

func newObjectLocationStore(coord gs.Coord, object obj.Objectable) *ObjectLocationStore {
	return &ObjectLocationStore{Coord: coord, ObjectId: object.ObjectId()}
}

func (store *ObjectLocationStore) Key() string {
	return fmt.Sprintf("%v:%v,%v", OBJECT_LOCATION_LAYER, store.Coord.X, store.Coord.Y)
}

func (store *ObjectLocationStore) Value() string {
	return store.ObjectId
}

type ObjectStore struct {
	ObjectId         string
	SerializedObject []byte
}

func newObjectStoreRetriever(objectId string) *ObjectStore {
	return &ObjectStore{ObjectId: objectId}
}

func newObjectStore(object obj.Objectable) *ObjectStore {
	serializedObject, _ := json.Marshal(object)

	return &ObjectStore{ObjectId: object.ObjectId(), SerializedObject: serializedObject}
}

func (store *ObjectStore) Key() string {
	return store.ObjectId
}

func (store *ObjectStore) Value() string {
	return string(store.SerializedObject)
}

type PaintLocationStore struct {
	Coord           gs.Coord
	SerializedPaint []byte
}

func newPaintLocationStore(coord gs.Coord, paint *pnt.Paint) *PaintLocationStore {
	serializedPaint, _ := json.Marshal(paint)
	return &PaintLocationStore{Coord: coord, SerializedPaint: serializedPaint}
}

func newPaintStoreRetriever(coord gs.Coord) *PaintLocationStore {
	return &PaintLocationStore{Coord: coord}
}

func (store *PaintLocationStore) Key() string {
	return fmt.Sprintf("%v:%v,%v", PAINT_LOCATION_LAYER, store.Coord.X, store.Coord.Y)
}

func (store *PaintLocationStore) Value() string {
	return string(store.SerializedPaint)
}
