package el

import (
	"../evt"
	"../items"
	"../obj"
	"../pnt"
	"../store"
	"../typ"
	"encoding/json"
	"fmt"
	"strings"
)

var EL_FACTORY *ElementFactory

type objFactory func() typ.Typical

type ElementFactory struct {
	factoryMap map[string]objFactory
}

func Factory() *ElementFactory {
	if EL_FACTORY == nil {
		fmt.Println("factory.go: No ElementFactory Instance found. Creating.")

		EL_FACTORY = &ElementFactory{
			factoryMap: make(map[string]objFactory),
		}

		EL_FACTORY.init()
	}
	return EL_FACTORY
}

func (elementFactory *ElementFactory) DeserializeObject(objectStore *store.ObjectStore) obj.Objectable {
	return elementFactory.deserialize(objectStore).(obj.Objectable)
}

func (elementFactory *ElementFactory) DeserializePaint(paintStore *store.PaintLocationStore) *pnt.Paint {
	return elementFactory.deserialize(paintStore).(*pnt.Paint)
}

func (elementFactory *ElementFactory) DeserializeItems(itemsStore *store.ItemsLocationStore) *items.Items {
	return elementFactory.deserialize(itemsStore).(*items.Items)
}

type actorDeserializer struct {
	Emitter store.TypeDeserializer
}

type actorDeserializerByte struct {
	Emitter string
}

func (elementFactory *ElementFactory) DeserializeEvent(serializedEvent string) *evt.Event {
	typical := elementFactory.load(evt.EVENT)
	json.Unmarshal([]byte(serializedEvent), &typical)
	event := typical.(*evt.Event)

	// TODO FIXME!!!!!!!

	fmt.Println(serializedEvent)
	emitterDeserializer := &actorDeserializer{}
	json.Unmarshal([]byte(serializedEvent), &emitterDeserializer)
	emitterType := emitterDeserializer.Emitter.Type
	fmt.Println(emitterType)

	actorDeserializerByte := &actorDeserializerByte{}
	json.Unmarshal([]byte(serializedEvent), &actorDeserializerByte)

	fmt.Println(actorDeserializerByte.Emitter)

	object := elementFactory.load(emitterType).(obj.Objectable)
	json.Unmarshal([]byte(actorDeserializerByte.Emitter), object)

	fmt.Println(object)

	event.Emitter = object

	return event
}

func (elementFactory *ElementFactory) deserialize(store store.Storable) typ.Typical {
	typical := elementFactory.load(store.GetType())
	json.Unmarshal(store.GetSerializedData(), &typical)

	return typical
}

func (elementFactory *ElementFactory) load(typeString string) typ.Typical {
	factory, ok := elementFactory.factoryMap[typeString]

	if !ok {
		// Factory has not been registered.
		// Make a list of all available factories and panic
		availableFactories := make([]string, len(elementFactory.factoryMap))
		for k := range elementFactory.factoryMap {
			availableFactories = append(availableFactories, k)
		}
		panic(fmt.Sprintf("Invalid Factory name. Must be one of: %s", strings.Join(availableFactories, ", ")))
	}

	return factory()
}

func (elementFactory *ElementFactory) register(name string, factory objFactory) {
	if factory == nil {
		panic("Cannot have nil dbo!")
	}

	_, registered := elementFactory.factoryMap[name]

	if registered {
		fmt.Printf("Objectable: %s already registered. Ignoring.\n", name)
	}

	elementFactory.factoryMap[name] = factory
}

func (elementFactory *ElementFactory) init() {
	elementFactory.register(obj.PLAYER, obj.LoadPlayer)
	elementFactory.register(pnt.PAINT, pnt.LoadPaint)
	elementFactory.register(items.ITEMS, items.LoadItems)
	elementFactory.register(evt.EVENT, evt.LoadEvent)
	fmt.Println("factory.go: Finished registering DboFactoring.")
}
