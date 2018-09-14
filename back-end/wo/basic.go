package wo

import (
	"../gs"
	"../items"
	"../terr"
	"math/rand"
)

//func storeElement(coord gs.Coord, dbo rc.Objectable) {
//	element := elementFactory.CreateNew(el.ELEMENT)
//	element.(*el.Element).DboKey = dbo.Key()
//	element.(*el.Element).GridCoord = coord
//
//	elementFactory.Save(element)
//}
//
//func LoadPlayer(id int) *el.Player {
//	return elementFactory.LoadFromId(el.PLAYER, id).(*el.Player)
//}

// creates a player
// returns a pair of Coord of World, SubWorld
//func InitializePlayer() *el.Player {
//	var player *el.Player
//
//	player = elementFactory.CreateNew(el.PLAYER).(*el.Player)
//	player.CoinCount = 0
//	player.Alive = true
//	player.Hp = 10
//
//	gridCoord := placeRandomLocation(player)
//
//	player.GridCoord = gridCoord
//
//	elementFactory.Save(player)
//
//	return player
//}

// creates a Rock
//func placeRockRandomly() {
//	rock := elementFactory.CreateNew(el.ROCK)
//	elementFactory.Save(rock)
//	placeRandomLocation(rock)
//}

//func MovePlayer(player *el.Player, vector gs.Vector) {
//	player.GridCoord = moveCharacter(player.GridCoord, vector, player)
//}
//
//func moveCharacter(coord gs.Coord, vector gs.Vector, element rc.Objectable) gs.Coord {
//	nextCoord, _ := SafeMove(coord, vector)
//
//	nextElement, _ := NextElement(coord, vector)
//
//	//fmt.Printf("basic.go: NextElement key: %s, value: %s\n", NextElement.String(), NextElement.Serialize())
//
//	override := false
//	switch element.(type) {
//	//case *el.Snake:
//	//	snake := element.(*el.Snake)
//	//	override = snake.Interact(NextElement)
//	case *el.Player:
//		player := element.(*el.Player)
//		override = player.Interact(nextElement)
//	default:
//		panic("I don't know how to move this type")
//	}
//
//	if override {
//		removeCoords(coord)
//		storeElement(nextCoord, element)
//	} else {
//		nextCoord = coord
//	}
//	return nextCoord
//}
//
//func elementFromKey(key string) rc.Objectable {
//	t, _ := rc.SplitKey(key)
//	element := elementFactory.LoadFromKey(t, key)
//
//	return element
//}
//
//func removeCoords(coord gs.Coord) {
//	location := el.NewLocation(coord)
//	element := elementFactory.LoadFromKey(el.ELEMENT, location.LocationKey()).(*el.Element)
//	elementFactory.Delete(element)
//}
//
//func elementFromCoords(coord gs.Coord) rc.Objectable {
//	location := el.NewLocation(coord)
//	element := elementFactory.LoadFromKey(el.ELEMENT, location.LocationKey()).(*el.Element)
//
//	if element.IsEmpty() {
//		//fmt.Printf("basic.go: Element is empty.\n")
//		return &el.Empty{}
//	} else {
//		//fmt.Printf("basic.go: Load From Key: %s\n", element.DboKey)
//		return elementFromKey(element.DboKey)
//	}
//}
//
//func NextElement(coord gs.Coord, vector gs.Vector) (rc.Objectable, bool) {
//	nextCoord, moved := SafeMove(coord, vector)
//	return elementFromCoords(nextCoord), moved
//}
//
//func SafeMove(gridCoord gs.Coord, vector gs.Vector) (gs.Coord, bool) {
//	gX := gridCoord.X + vector.X
//	gY := gridCoord.Y + vector.Y
//
//	if isOutOfBound(gX, gY, gs.WORLD_SIZE) {
//		return gridCoord, false
//	}
//
//	return gs.NewCoord(gX, gY), true
//}

func InitializeWorld() {
	for i := 0; i < gs.WORLD_SIZE; i++ {
		for j := 0; j < gs.WORLD_SIZE; j++ {
			coord := gs.NewCoord(i, j)
			coinChance := rand.Intn(10)
			rockChance := rand.Intn(8)

			if coinChance == 0 {
				coins := items.NewCoinStack()
				AddItemsToStack(coord, coins)
			} else if rockChance == 0 {
				SetTerrain(coord, terr.NewRock())
			}
		}
	}
}
