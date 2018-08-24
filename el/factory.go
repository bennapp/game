package el

import (
	"../rc"
	"fmt"
	"strconv"
	"strings"
)

const GLOBALID = "globalId"

var INSTANCE *ElementFactory

type ElementFactory struct {
	dboManager *rc.RedisManager
	factoryMap map[string]DboFactory
	globalId   int
}

func Factory(firstBoot bool) *ElementFactory {
	if INSTANCE == nil {
		fmt.Println("factory.go: No ElementFactory Instance found. Creating.")

		dboManager := rc.Manager()
		globalId := 1

		if !firstBoot {
			globalId, _ = strconv.Atoi(dboManager.LoadFromKey(GLOBALID))
		}

		INSTANCE = &ElementFactory{
			dboManager: dboManager,
			factoryMap: make(map[string]DboFactory),
			globalId:   globalId,
		}
	}

	return INSTANCE
}

//creates a blank Dbo with a new Id
func (elementFactory *ElementFactory) CreateNew(t string) rc.Dbo {
	return elementFactory.Create(t, true)
}

//creates a Dbo
func (elementFactory *ElementFactory) Create(t string, isNew bool) rc.Dbo {
	factory, ok := elementFactory.factoryMap[t]

	if !ok {
		// Factory has not been registered.
		// Make a list of all available factories and panic
		availableFactories := make([]string, len(elementFactory.factoryMap))
		for k, _ := range elementFactory.factoryMap {
			availableFactories = append(availableFactories, k)
		}
		panic(fmt.Sprintf("Invalid Factory name. Must be one of: %s", strings.Join(availableFactories, ", ")))
	}

	if isNew {
		return factory(elementFactory.nextGlobalId())
	} else {
		return factory(-1)
	}
}

func (elementFactory *ElementFactory) LoadFromId(t string, id int) rc.Dbo {
	return elementFactory.LoadFromKey(t, rc.GenerateKey(t, id))
}

func (elementFactory *ElementFactory) LoadFromKey(t string, key string) rc.Dbo {
	blankDbo := elementFactory.Create(t, false)

	val := elementFactory.dboManager.LoadFromKey(key)

	blankDbo.Deserialize(key, val)

	return blankDbo
}

func (elementFactory *ElementFactory) Delete(dbo rc.Dbo) {
	elementFactory.dboManager.Delete(dbo)
}

func (elementFactory *ElementFactory) Save(dbo rc.Dbo) {
	elementFactory.dboManager.Save(dbo)
}

type DboFactory func(id int) rc.Dbo

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

func (elementFactory *ElementFactory) Init() {
	elementFactory.Register(COIN, newCoinDbo)
	elementFactory.Register(ROCK, newRockDbo)
	elementFactory.Register(PLAYER, newPlayerDbo)
	elementFactory.Register(ELEMENT, newElementDbo)
	fmt.Println("factory.go: Finished Factory init.")
}

func (elementFactory *ElementFactory) nextGlobalId() int {
	elementFactory.globalId++
	return elementFactory.globalId
}
