package wo

import (
	"../el"
	"../gs"
	"../rc"
	"fmt"
	"math/rand"
	"time"
)

const MAX_COIN_AMOUNT = 10

// GLOBALS
var elementFactory *el.ElementFactory

func isEmpty(coord gs.Coord) bool {
	// Look into this, not sure this is right
	// can we check that value == nil (?)
	location := el.NewLocation(coord)
	element := elementFactory.LoadFromKey(el.ELEMENT, location.LocationKey())

	return element.(*el.Element).IsEmpty()
}

func storeElement(coord gs.Coord, dbo rc.Dbo) {
	element := elementFactory.CreateNew(el.ELEMENT)
	element.(*el.Element).DboKey = dbo.Key()
	element.(*el.Element).GridCoord = coord

	elementFactory.Save(element)
}

func LoadPlayer(id int) *el.Player {
	return elementFactory.LoadFromId(el.PLAYER, id).(*el.Player)
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

func buildAndStoreCoin() *el.Coin {
	coin := initializeCoin()

	elementFactory.Save(coin)

	return coin
}

// creates a Coin
func placeCoin(coin *el.Coin) {
	placeElementRandomLocationWithLock(coin)
}

// creates a Rock
func placeRock() {
	rock := elementFactory.CreateNew(el.ROCK)
	elementFactory.Save(rock)
	placeElementRandomLocationWithLock(rock)
}

func placeElementRandomLocationWithLock(dbo rc.Dbo) gs.Coord {
	coord := randomCoord()

	if isEmpty(coord) {
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
	nextCoord := gs.NewCoord(coord.X+vector.X, coord.Y+vector.Y)

	if isOutOfBound(nextCoord.X, nextCoord.Y, gs.GRID_SIZE) {
		return coord
	}

	nextElement, _ := NextElement(coord, vector)

	//fmt.Printf("basic.go: NextElement key: %s, value: %s\n", NextElement.String(), NextElement.Serialize())

	override := false
	switch element.(type) {
	//case *el.Snake:
	//	snake := element.(*el.Snake)
	//	override = snake.Interact(NextElement)
	case *el.Player:
		player := element.(*el.Player)
		override = player.Interact(nextElement)
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
	t, _ := rc.SplitKey(key)
	element := elementFactory.LoadFromKey(t, key)

	return element
}

//TODO - create removeCoords in manager.go
func removeCoords(coord gs.Coord) {
	location := el.NewLocation(coord)
	element := elementFactory.LoadFromKey(el.ELEMENT, location.LocationKey()).(*el.Element)
	elementFactory.Delete(element)
}

func elementFromCoords(coord gs.Coord) rc.Dbo {
	location := el.NewLocation(coord)
	element := elementFactory.LoadFromKey(el.ELEMENT, location.LocationKey()).(*el.Element)

	if element.IsEmpty() {
		//fmt.Printf("basic.go: Element is empty.\n")
		return &el.Empty{}
	} else {
		//fmt.Printf("basic.go: Load From Key: %s\n", element.DboKey)
		return elementFromKey(element.DboKey)
	}
}

func NextElement(coord gs.Coord, vector gs.Vector) (rc.Dbo, bool) {
	nextCoord := gs.NewCoord(coord.X+vector.X, coord.Y+vector.Y)

	if isOutOfBound(nextCoord.X, nextCoord.Y, gs.GRID_SIZE) {
		return &el.Empty{}, false
	}

	return elementFromCoords(nextCoord), true
}

func subWorldMove(subWorldCoord gs.Coord, gridCoord gs.Coord, vector gs.Vector) (gs.Coord, gs.Coord, bool) {
	wX := subWorldCoord.X + carry(gridCoord.X, vector.X, gs.GRID_SIZE)
	wY := subWorldCoord.Y + carry(gridCoord.Y, vector.Y, gs.GRID_SIZE)

	gX := wrap(gridCoord.X, vector.X, gs.GRID_SIZE)
	gY := wrap(gridCoord.Y, vector.Y, gs.GRID_SIZE)

	if isOutOfBound(wX, wY, gs.WORLD_SIZE) {
		return subWorldCoord, gridCoord, false
	}

	return gs.NewCoord(wX, wY), gs.NewCoord(gX, gY), true
}

func SpawnCoinsInWorld() {
	sleepTime := 10000 * time.Millisecond

	for {
		fmt.Printf("basic.go: spawning random coin.\n")
		spawnCoinInWorld()
		time.Sleep(sleepTime)
		sleepTime += sleepTime
	}
}

func spawnCoinInWorld() {
	coin := buildAndStoreCoin()
	placeCoin(coin)
}

func InitializeWorld() {
	for i := 0; i < gs.GRID_SIZE*2; i++ {
		placeRock()
	}
}

func Init() {
	elementFactory = el.Factory()
}

func Reset() {
	elementFactory.Reset()
}

func Close() {
	elementFactory.Close()
}

func PlayerLogout(player *el.Player) {
	elementFactory.Save(player)
}
