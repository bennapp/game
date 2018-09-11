package store

import (
	"../gs"
	"../items"
	"../obj"
	"../pnt"
	"encoding/json"
	"fmt"
)

const OBJECT_LOCATION_LAYER = "o"
const PAINT_LOCATION_LAYER = "p"
const ITEMS_LOCATION_LAYER = "i"

type RedisStore interface {
	Key() string
	Value() string
}

type Storable interface {
	GetType() string
	GetSerializedData() []byte
}

type ObjectLocationStore struct {
	Coord    gs.Coord
	ObjectId string
}

func NewObjectLocationStoreRetriever(coord gs.Coord) *ObjectLocationStore {
	return &ObjectLocationStore{Coord: coord}
}

func NewObjectLocationStore(coord gs.Coord, object obj.Objectable) *ObjectLocationStore {
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
	Type             string
}

func NewObjectStoreRetriever(objectId string) *ObjectStore {
	return &ObjectStore{ObjectId: objectId}
}

func NewObjectStore(object obj.Objectable) *ObjectStore {
	serializedObject, _ := json.Marshal(object)

	return &ObjectStore{ObjectId: object.ObjectId(), SerializedObject: serializedObject}
}

func (store *ObjectStore) Key() string {
	return store.ObjectId
}

func (store *ObjectStore) Value() string {
	return string(store.SerializedObject)
}

func (store *ObjectStore) GetType() string {
	return store.Type
}

func (store *ObjectStore) GetSerializedData() []byte {
	return store.SerializedObject
}

func (store *ObjectStore) Retrieve(objectData string) {
	store.SerializedObject = []byte(objectData)
	store.Type = newTypeDeserializer(store.SerializedObject).Type
}

type PaintLocationStore struct {
	Coord           gs.Coord
	SerializedPaint []byte
}

func NewPaintLocationStore(coord gs.Coord, paint *pnt.Paint) *PaintLocationStore {
	serializedPaint, _ := json.Marshal(paint)
	return &PaintLocationStore{Coord: coord, SerializedPaint: serializedPaint}
}

func NewPaintStoreRetriever(coord gs.Coord) *PaintLocationStore {
	return &PaintLocationStore{Coord: coord}
}

func (store *PaintLocationStore) Key() string {
	return fmt.Sprintf("%v:%v,%v", PAINT_LOCATION_LAYER, store.Coord.X, store.Coord.Y)
}

func (store *PaintLocationStore) Value() string {
	return string(store.SerializedPaint)
}

func (store *PaintLocationStore) GetType() string {
	return pnt.PAINT
}

func (store *PaintLocationStore) GetSerializedData() []byte {
	return store.SerializedPaint
}

type ItemsLocationStore struct {
	Coord           gs.Coord
	SerializedItems []byte
}

func NewItemsLocationStore(coord gs.Coord, items *items.Items) *ItemsLocationStore {
	serializedItems, _ := json.Marshal(items)
	return &ItemsLocationStore{Coord: coord, SerializedItems: serializedItems}
}

func NewItemsStoreRetriever(coord gs.Coord) *ItemsLocationStore {
	return &ItemsLocationStore{Coord: coord}
}

func (store *ItemsLocationStore) Key() string {
	return fmt.Sprintf("%v:%v,%v", ITEMS_LOCATION_LAYER, store.Coord.X, store.Coord.Y)
}

func (store *ItemsLocationStore) Value() string {
	return string(store.SerializedItems)
}

func (store *ItemsLocationStore) GetType() string {
	return items.ITEMS
}

func (store *ItemsLocationStore) GetSerializedData() []byte {
	return store.SerializedItems
}

type TypeDeserializer struct {
	Type string
}

func newTypeDeserializer(serializedData []byte) *TypeDeserializer {
	typeDeserializer := new(TypeDeserializer)

	json.Unmarshal(serializedData, typeDeserializer)

	return typeDeserializer
}
