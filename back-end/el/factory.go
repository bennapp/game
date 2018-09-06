package el

import (
	"../rc"
	"encoding/json"
	"fmt"
		"strings"
	"github.com/google/uuid"
)

var INSTANCE *ElementFactory

type ElementFactory struct {
	dboManager *rc.RedisManager
	factoryMap map[string]DboFactory
	globalId   int
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
	return elementFactory.Create(t)
}

//creates a Dbo
func (elementFactory *ElementFactory) Create(t string) rc.Dbo {
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

	return factory(uuid.New())
}

func (elementFactory *ElementFactory) LoadFromId(elType string, id uuid.UUID) rc.Dbo {
	return elementFactory.LoadFromKey(elType, rc.GenerateKey(id))
}

func (elementFactory *ElementFactory) LoadFromKey(elType string, key string) rc.Dbo {
	blankDbo := elementFactory.Load(elType)

	val := []byte(elementFactory.dboManager.LoadFromKey(key))

	json.Unmarshal(val, &blankDbo)

	return blankDbo
}

func (elementFactory *ElementFactory) Delete(dbo rc.Dbo) {
	elementFactory.dboManager.Delete(dbo)
}

func (elementFactory *ElementFactory) Save(dbo rc.Dbo) {
	elementFactory.dboManager.Save(dbo)
}

type DboFactory func(id uuid.UUID) rc.Dbo

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
	fmt.Println("factory.go: Finished registering DboFactoring.")
}

func (elementFactory *ElementFactory) Reset() {
	elementFactory.dboManager.Client("ALL").FlushAll()
}
