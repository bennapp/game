package el

import (
	"../rc"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"strings"
)

var INSTANCE *ElementFactory

type ElementFactory struct {
	dboManager *rc.RedisManager
	factoryMap map[string]DboFactory
}

func Factory() *ElementFactory {
	if INSTANCE == nil {
		fmt.Println("factory.go: No ElementFactory Instance found. Creating.")

		dboManager := rc.Manager()

		INSTANCE = &ElementFactory{
			dboManager: dboManager,
			factoryMap: make(map[string]DboFactory),
		}

		INSTANCE.init()
	}
	return INSTANCE
}

//creates a blank Dbo with a new id
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

	return factory(isNew)
}

// can this use LoadFromKeyWithoutType ?
func (elementFactory *ElementFactory) LoadFromId(elType string, id uuid.UUID) rc.Dbo {
	return elementFactory.LoadFromKey(elType, rc.GenerateKey(id))
}

func (elementFactory *ElementFactory) LoadFromKey(elType string, key string) rc.Dbo {
	blankDbo := elementFactory.Create(elType, false)

	values, empty := elementFactory.dboManager.LoadFromKey(key)

	if !empty {
		json.Unmarshal([]byte(values), &blankDbo)

		return blankDbo
	} else {
		return nil
	}
}

func (elementFactory *ElementFactory) Delete(dbo rc.Dbo) {
	elementFactory.dboManager.Delete(dbo)
}

func (elementFactory *ElementFactory) Save(dbo rc.Dbo) {
	elementFactory.dboManager.Save(dbo)
}

type DboFactory func(bool) rc.Dbo

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
	elementFactory.Register(COIN, newCoinDbo)
	elementFactory.Register(PLAYER, newPlayerDbo)
	elementFactory.Register(ELEMENT, newElementDbo)
	fmt.Println("factory.go: Finished registering DboFactoring.")
}

func (elementFactory *ElementFactory) Reset() {
	elementFactory.dboManager.Client("ALL").FlushAll()
}
