package el

import (
	"../rc"
	"fmt"
	"strings"
)

var EL_FACTORY *ElementFactory

type ElementFactory struct {
	dboManager *rc.RedisManager
	factoryMap map[string]DboFactory
}

func Factory() *ElementFactory {
	if EL_FACTORY == nil {
		fmt.Println("factory.go: No ElementFactory Instance found. Creating.")

		dboManager := rc.Manager()

		EL_FACTORY = &ElementFactory{
			dboManager: dboManager,
			factoryMap: make(map[string]DboFactory),
		}

		EL_FACTORY.init()
	}
	return EL_FACTORY
}

func (elementFactory *ElementFactory) Load(objectType string) rc.Dbo {
	factory, ok := elementFactory.factoryMap[objectType]

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

type DboFactory func() rc.Dbo

func (elementFactory *ElementFactory) Register(name string, factory DboFactory) {
	if factory == nil {
		panic("Cannot have nil dbo!")
	}

	_, registered := elementFactory.factoryMap[name]

	if registered {
		fmt.Printf("Dbo: %s already registered. Ignoring.\n", name)
	}

	elementFactory.factoryMap[name] = factory
}

func (elementFactory *ElementFactory) init() {
	elementFactory.Register(COIN, loadCoin)
	//elementFactory.Register(ROCK, newRockDbo)
	//elementFactory.Register(PLAYER, newPlayerDbo)
	//elementFactory.Register(ELEMENT, newElementDbo)
	fmt.Println("factory.go: Finished registering DboFactoring.")
}
