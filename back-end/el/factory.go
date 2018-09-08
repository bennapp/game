package el

import (
	"../gs"
	"../obj"
	"../rc"
	"../pnt"
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

func (elementFactory *ElementFactory) Load(typeString string) typ.Typical {
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

func (elementFactory *ElementFactory) Register(name string, factory objFactory) {
	if factory == nil {
		panic("Cannot have nil dbo!")
	}

	_, registered := elementFactory.factoryMap[name]

	if registered {
		fmt.Printf("Objectable: %s already registered. Ignoring.\n", name)
	}

	elementFactory.factoryMap[name] = factory
}

func (elementFactory *ElementFactory) LoadObjectFromCoord(coord gs.Coord) typ.Typical {
	objectType, objectStore := rc.Manager().LoadObjectTypeFromCoord(coord)

	if objectType == "" {
		return nil
	}

	object := elementFactory.Load(objectType)
	json.Unmarshal(objectStore.SerializedObject, &object)

	return object
}

func (elementFactory *ElementFactory) LoadPaintFromCoord(coord gs.Coord) typ.Typical {
	paintStore := rc.Manager().LoadPaintStoreFromCoord(coord)

	if paintStore == nil {
		return nil
	}

	paint := elementFactory.Load(pnt.PAINT)
	json.Unmarshal(paintStore.SerializedPaint, &paint)

	return paint
}

func (elementFactory *ElementFactory) init() {
	elementFactory.Register(obj.COIN, obj.LoadCoin)
	elementFactory.Register(pnt.PAINT, pnt.LoadPaint)
	//elementFactory.Register(terr.ROCK, terr.LoadRock)
	//elementFactory.Register(PLAYER, newPlayerDbo)
	//elementFactory.Register(ELEMENT, newElementDbo)
	fmt.Println("factory.go: Finished registering DboFactoring.")
}
