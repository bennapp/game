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

func isEmpty(subWorldCoord gs.Coord, coord gs.Coord) bool {
	// Look into this, not sure this is right
	// can we check that value == nil (?)
	location := el.NewLocation(subWorldCoord, coord)
	element := elementFactory.LoadFromKey(el.ELEMENT, location.LocationKey())

	return element.(*el.Element).IsEmpty()
}

func storeElement(subWorldCoord gs.Coord, coord gs.Coord, dbo rc.Dbo) {
	element := elementFactory.CreateNew(el.ELEMENT)
	element.(*el.Element).DboKey = dbo.Key()
	element.(*el.Element).SubWorldCoord = subWorldCoord
	element.(*el.Element).GridCoord = coord

	elementFactory.Save(element)
}

func LoadPlayer(id int) *el.Player {
	return elementFactory.LoadFromId(el.PLAYER, id).(*el.Player)
}

// creates a player
// returns a pair of Coord of World, SubWorld
func InitializePlayer() *el.Player {
	x, y := randomPair(gs.WORLD_SIZE)

	var player *el.Player

	player = elementFactory.CreateNew(el.PLAYER).(*el.Player)
	player.CoinCount = 0
	player.Alive = true
	player.Hp = 10

	subWorldCoord := gs.NewCoord(x, y)
	gridCoord := placeElementRandomLocationWithLock(subWorldCoord, player)

	player.SubWorldCoord = subWorldCoord
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
func placeCoin(subWorldCoord gs.Coord, coin *el.Coin) {
	placeElementRandomLocationWithLock(subWorldCoord, coin)
}

// creates a Rock
func placeRock(subWorldCoord gs.Coord) {
	rock := elementFactory.CreateNew(el.ROCK)
	elementFactory.Save(rock)
	placeElementRandomLocationWithLock(subWorldCoord, rock)
}

func placeElementRandomLocationWithLock(subWorldCoord gs.Coord, dbo rc.Dbo) gs.Coord {
	//TODO - refactor this to a method in coord
	x, y := randomPair(gs.GRID_SIZE)
	coord := gs.NewCoord(x, y)

	if isEmpty(subWorldCoord, coord) {
		storeElement(subWorldCoord, coord, dbo)
	} else {
		coord = placeElementRandomLocationWithLock(subWorldCoord, dbo)
	}

	return coord
}

func MovePlayer(player *el.Player, vector gs.Vector) {
	player.SubWorldCoord, player.GridCoord = moveCharacter(player.SubWorldCoord, player.GridCoord, vector, player)
}

func moveCharacter(subWorldCoord gs.Coord, coord gs.Coord, vector gs.Vector, element rc.Dbo) (gs.Coord, gs.Coord) {
	nextSubWorldCoord, nextCoord, _ := subWorldMove(subWorldCoord, coord, vector)

	nextElement, _ := NextElement(subWorldCoord, coord, vector)

	//fmt.Printf("game.go: NextElement key: %s, value: %s\n", NextElement.String(), NextElement.Serialize())

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
		removeCoords(subWorldCoord, coord)
		storeElement(nextSubWorldCoord, nextCoord, element)
	} else {
		nextSubWorldCoord = subWorldCoord
		nextCoord = coord
	}
	return nextSubWorldCoord, nextCoord
}

func elementFromKey(key string) rc.Dbo {
	t, _ := rc.SplitKey(key)
	element := elementFactory.LoadFromKey(t, key)

	return element
}

//TODO - create removeCoords in manager.go
func removeCoords(subWorldCoord gs.Coord, coord gs.Coord) {
	location := el.NewLocation(subWorldCoord, coord)
	element := elementFactory.LoadFromKey(el.ELEMENT, location.LocationKey()).(*el.Element)
	elementFactory.Delete(element)
}

func elementFromCoords(subWorldCoord gs.Coord, coord gs.Coord) rc.Dbo {
	location := el.NewLocation(subWorldCoord, coord)
	element := elementFactory.LoadFromKey(el.ELEMENT, location.LocationKey()).(*el.Element)

	if element.IsEmpty() {
		//fmt.Printf("game.go: Element is empty.\n")
		return &el.Empty{}
	} else {
		//fmt.Printf("game.go: Load From Key: %s\n", element.DboKey)
		return elementFromKey(element.DboKey)
	}
}

func NextElement(subWorldCoord gs.Coord, coord gs.Coord, vector gs.Vector) (rc.Dbo, bool) {
	nextSubWorldCoord, nextCoord, moved := subWorldMove(subWorldCoord, coord, vector)
	return elementFromCoords(nextSubWorldCoord, nextCoord), moved
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
		randomSubWorldCoord := randomSubWorldCoord()
		spawnCoinInSubWorld(randomSubWorldCoord)
		time.Sleep(sleepTime)
		sleepTime += sleepTime
	}
}

func spawnCoinInSubWorld(subWorldCoord gs.Coord) {
	coin := buildAndStoreCoin()
	placeCoin(subWorldCoord, coin)
}

func InitializeWorld() {
	for i := 0; i < gs.WORLD_SIZE; i++ {
		for j := 0; j < gs.WORLD_SIZE; j++ {
			subWorldCoord := gs.NewCoord(i, j)
			initializeSubWorld(subWorldCoord)
		}
	}
}

func initializeSubWorld(subWorldCoord gs.Coord) {
	//TODO - remove 10/2 hard coded value
	for i := 0; i < 10; i++ {
		placeRock(subWorldCoord)
	}

	for i := 0; i < 2; i++ {
		coin := buildAndStoreCoin()
		placeCoin(subWorldCoord, coin)
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
