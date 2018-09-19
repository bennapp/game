package typf

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

var TYPE_FACTORY *TypeFactory

type objFactory func() typ.Typical

type TypeFactory struct {
	factoryMap map[string]objFactory
}

func Factory() *TypeFactory {
	if TYPE_FACTORY == nil {
		fmt.Println("factory.go: No TypeFactory Instance found. Creating.")

		TYPE_FACTORY = &TypeFactory{
			factoryMap: make(map[string]objFactory),
		}

		TYPE_FACTORY.init()
	}
	return TYPE_FACTORY
}

func (elementFactory *TypeFactory) DeserializeObject(objectStore *store.ObjectStore) obj.Objectable {
	typical := elementFactory.deserialize(objectStore)

	if typical == nil {
		return nil
	}

	return typical.(obj.Objectable)
}

func (elementFactory *TypeFactory) DeserializePaint(paintStore *store.PaintLocationStore) *pnt.Paint {
	typical := elementFactory.deserialize(paintStore)

	if typical == nil {
		return nil
	}

	return typical.(*pnt.Paint)
}

func (elementFactory *TypeFactory) DeserializeItems(itemsStore *store.ItemsLocationStore) *items.Items {
	typical := elementFactory.deserialize(itemsStore)

	if typical == nil {
		return nil
	}

	return typical.(*items.Items)
}

func (elementFactory *TypeFactory) DeserializeEvent(serializedEvent string) *evt.Event {
	typical := elementFactory.load(evt.EVENT)
	json.Unmarshal([]byte(serializedEvent), &typical)
	event := typical.(*evt.Event)

	return event
}

func (elementFactory *TypeFactory) deserialize(store store.DeserializableStorable) typ.Typical {
	typical := elementFactory.load(store.GetType())
	json.Unmarshal(store.GetSerializedData(), &typical)

	return typical
}

func (elementFactory *TypeFactory) load(typeString string) typ.Typical {
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

func (elementFactory *TypeFactory) register(name string, factory objFactory) {
	if factory == nil {
		panic("Cannot have nil dbo!")
	}

	_, registered := elementFactory.factoryMap[name]

	if registered {
		fmt.Printf("Objectable: %s already registered. Ignoring.\n", name)
	}

	elementFactory.factoryMap[name] = factory
}

func (elementFactory *TypeFactory) init() {
	elementFactory.register(obj.PLAYER, obj.LoadPlayer)
	elementFactory.register(pnt.PAINT, pnt.LoadPaint)
	elementFactory.register(items.ITEMS, items.LoadItems)
	elementFactory.register(evt.EVENT, evt.LoadEvent)
	fmt.Println("factory.go: Finished registering DboFactoring.")
}
