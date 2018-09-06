package wo

import (
	"../el"
	"../gs"
	"../rc"
	"github.com/google/uuid"
	"math/rand"
)

const MAX_COIN_AMOUNT = 10

// GLOBALS
var elementFactory *el.ElementFactory

func IsEmpty(coord gs.Coord) bool {
	location := el.NewLocation(coord)
	element := elementFactory.LoadFromKeyWithoutType(location.LocationKey())

	return element.(*el.Element).IsEmpty()
}

func storeElement(coord gs.Coord, dbo rc.Dbo) {
	element := elementFactory.CreateNew(el.ELEMENT)
	element.(*el.Element).DboKey = dbo.Key()
	element.(*el.Element).GridCoord = coord

	elementFactory.Save(element)
}

func LoadPlayer(id string) *el.Player {
	return elementFactory.LoadFromId(el.PLAYER, uuid.MustParse(id)).(*el.Player)
}

// creates a player
// returns a pair of Coord of World, SubWorld
func InitializePlayer() *el.Player {
	var player *el.Player

	player = elementFactory.CreateNew(el.PLAYER).(*el.Player)
	player.CoinCount = 0
	player.Alive = true
	player.Hp = 10

	gridCoord := placeElementRandomLocationWithLock(player)

	player.GridCoord = gridCoord

	elementFactory.Save(player)

	return player
}

func initializeCoin() *el.Coin {
	coin := elementFactory.CreateNew(el.COIN)
	coin.(*el.Coin).Amount = rand.Intn(MAX_COIN_AMOUNT)

	return coin.(*el.Coin)
}

func BuildAndStoreCoin() *el.Coin {
	coin := initializeCoin()

	elementFactory.Save(coin)

	return coin
}

// creates a Coin
func PlaceCoinRandomly(coin *el.Coin) {
	placeElementRandomLocationWithLock(coin)
}

// creates a Rock
func placeRockRandomly() {
	//rock := elementFactory.CreateNew(el.ROCK)
	//elementFactory.Save(rock)
	//placeElementRandomLocationWithLock(rock)
}

func placeElementRandomLocationWithLock(dbo rc.Dbo) gs.Coord {
	//TODO - refactor this to a method in coord
	x, y := randomPair(gs.WORLD_SIZE)
	coord := gs.NewCoord(x, y)

	if IsEmpty(coord) {
		storeElement(coord, dbo)
	} else {
		coord = placeElementRandomLocationWithLock(dbo)
	}

	return coord
}

func MovePlayer(player *el.Player, vector gs.Vector) {
	player.GridCoord = moveCharacter(player.GridCoord, vector, player)
}

func moveCharacter(coord gs.Coord, vector gs.Vector, element rc.Dbo) gs.Coord {
	nextCoord, _ := SafeMove(coord, vector)

	nextElement, _ := NextElement(coord, vector)

	//fmt.Printf("basic.go: NextElement key: %s, value: %s\n", NextElement.String(), NextElement.Serialize())

	override := false
	switch element.(type) {
	//case *el.Snake:
	//	snake := element.(*el.Snake)
	//	override = snake.Interact(NextElement)
	//	if override {
	//		elementFactory.Save(snake)
	//	}
	case *el.Player:
		player := element.(*el.Player)
		override = player.Interact(nextElement)

		if override {
			elementFactory.Save(player)
		}
	default:
		panic("I don't know how to move this type")
	}

	if override {
		removeCoords(coord)
		storeElement(nextCoord, element)
	} else {
		nextCoord = coord
	}
	return nextCoord
}

func elementFromKey(key string) rc.Dbo {
	element := elementFactory.LoadFromKeyWithoutType(key)

	return element
}

//TODO - create removeCoords in manager.go
func removeCoords(coord gs.Coord) {
	location := el.NewLocation(coord)
	element := elementFactory.LoadFromKeyWithoutType(location.LocationKey()).(*el.Element)
	elementFactory.Delete(element)
}

func elementFromCoords(coord gs.Coord) rc.Dbo {
	location := el.NewLocation(coord)
	element := elementFactory.LoadFromKeyWithoutType(location.LocationKey()).(*el.Element)

	if element.IsEmpty() {
		return &el.Empty{}
	} else {
		return elementFromKey(element.DboKey)
	}
}

func NextElement(coord gs.Coord, vector gs.Vector) (rc.Dbo, bool) {
	nextCoord, moved := SafeMove(coord, vector)
	return elementFromCoords(nextCoord), moved
}

func SafeMove(gridCoord gs.Coord, vector gs.Vector) (gs.Coord, bool) {
	gX := gridCoord.X + vector.X
	gY := gridCoord.Y + vector.Y

	if isOutOfBound(gX, gY, gs.WORLD_SIZE) {
		return gridCoord, false
	}

	return gs.NewCoord(gX, gY), true
}

func InitializeWorld() {
	for i := 0; i < gs.WORLD_SIZE; i++ {
		for j := 0; j < gs.WORLD_SIZE; j++ {
			rockChance := rand.Intn(5)
			if rockChance == 0 {
				placeRockRandomly()
			}

			coinChance := rand.Intn(10)
			if coinChance == 0 {
				coin := BuildAndStoreCoin()
				PlaceCoinRandomly(coin)
			}
		}
	}
}

func Init() {
	elementFactory = el.Factory()
}

func Reset() {
	elementFactory.Reset()
}
