package el

import (
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
		for k, _ := range elementFactory.factoryMap {
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
	elementFactory.register(obj.COIN, obj.LoadCoin)
	elementFactory.register(pnt.PAINT, pnt.LoadPaint)
	//elementFactory.Register(terr.ROCK, terr.LoadRock)
	//elementFactory.Register(PLAYER, newPlayerDbo)
	//elementFactory.Register(ELEMENT, newElementDbo)
	fmt.Println("factory.go: Finished registering DboFactoring.")
}
