package main

/*
import "github.com/nsf/termbox-go"
import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"
)
import (
	"./el"
	"./gs"
	"./rc"
)

const MAX_COIN_AMOUNT = 10

// GLOBALS
var elementFactory *el.ElementFactory
var firstBoot bool
var debug bool
var coinIdInc int

//func redisSet(key string, value string) {
//	err := redisClient.Set(key, value, 0).Err()
//	if err != nil {
//		panic(err)
//	}
//}
//func storeObject(id string, object string) {
//	redisSet(id, object)
//}
//
//func storePlayer(player *el.Player) {
//	storeObject(player.id(), player.Val())
//}
//func storeCoin(coin *el.Coin) {
//	storeObject(coin.id(), coin.Val())
//}
//
//func coordKey(subWorldCoord gs.Coord, coord gs.Coord) string {
//	return fmt.Sprintf("%v:%v", subWorldCoord.Key(), coord.Key())
//}
//func subWorldCoordKey(subWorld *gs.SubWorld, coord gs.Coord) string {
//	return coordKey(subWorld.Coord(), coord)
//}
//
//func storeCoord(subWorld *gs.SubWorld, coord gs.Coord, id string) {
//	redisSet(subWorldCoordKey(subWorld, coord), id)
//}
//
//func storeCoinCoord(subWorld *gs.SubWorld, coord gs.Coord, coin *el.Coin) {
//	storeCoord(subWorld, coord, coin.id())
//}
//func storeRockCoord(subWorld *gs.SubWorld, coord gs.Coord, rock *el.Rock) {
//	storeCoord(subWorld, coord, rock.id())
//}
//func storePlayerCoord(subWorld *gs.SubWorld, coord gs.Coord, player *el.Player) {
//	storeCoord(subWorld, coord, player.id())
//}
//
//func emptykey(key string) {
//	err := redisClient.Set(key, nil, 0).Err()
//	if err != nil {
//		panic(err)
//	}
//}
//
//func setEmptyValue(subWorldCoord gs.Coord, coord gs.Coord) {
//	emptykey(coordKey(subWorldCoord, coord))
//}
//
//func setEmptyObject(id string) {
//	emptykey(id)
//}
//
func isEmpty(subWorldCoord gs.Coord, coord gs.Coord) bool {
	// Look into this, not sure this is right
	// can we check that value == nil (?)
	location := el.NewLocation(subWorldCoord, coord)
	element := elementFactory.LoadFromKey(el.ELEMENT, location.LocationKey())

	return element.(*el.Element).IsEmpty()
}

func storeElement(subWorld *gs.SubWorld, coord gs.Coord, dbo rc.Dbo) {
	element := elementFactory.CreateNew(el.ELEMENT)
	element.(*el.Element).DboKey = dbo.Key()
	element.(*el.Element).SubWorldCoord = subWorld.Coord()
	element.(*el.Element).GridCoord = coord

	elementFactory.Save(element)
}

//func (player *Player) BuildWall(gs *World) bool {
//	wallCost := 5
//
//	if player.coinCount < wallCost {
//		return false
//	}
//
//	up := Vector{x: 0, y: -1}
//	upRight := Vector{x: 1, y: -1}
//
//	upCell := nextCell(gs, player.SubWorldCoord, player.GridCoord, up)
//	upElement := upCell.element
//
//	upRightCell := nextCell(gs, player.SubWorldCoord, player.GridCoord, upRight)
//	upRightElement := upRightCell.element
//
//	if isEmpty(upElement) && isEmpty(upRightElement) {
//		player.mux.Lock()
//		player.coinCount -= wallCost
//		player.mux.Unlock()
//
//		upCell.mux.Lock()
//		upCell.element = &Building{code: "<"}
//		upCell.mux.Unlock()
//
//		upRightCell.mux.Lock()
//		upRightCell.element = &Building{code: ">"}
//		upRightCell.mux.Unlock()
//
//		return true
//	}
//
//	return false
//}
func playerKey(id int) string {
	return fmt.Sprintf("player:%v", id)
}

// creates a player
// returns a pair of Coord of World, SubWorld
func initializePlayer(wo *gs.World) *el.Player {
	x, y := randomPair(gs.WORLD_SIZE)
	subWorld := wo.SubWorlds()[x][y]

	var player *el.Player
	if firstBoot {
		player = elementFactory.CreateNew(el.PLAYER).(*el.Player)
		player.CoinCount = 0
		player.Alive = true
		player.Hp = 10
		player.id = 1 //TODO - hard coded id; fix??

		subWorldCoord := gs.NewCoord(x, y)
		gridCoord := placeElementRandomLocationWithLock(&subWorld, player)

		player.SubWorldCoord = subWorldCoord
		player.GridCoord = gridCoord

		elementFactory.Save(player)
	} else {
		player = elementFactory.LoadFromId(el.PLAYER, 1).(*el.Player)
	}

	return player
}

// creates a snakeCell
// returns a pair of Coord of World, SubWorld
//func initializeSnake(wo *gs.World) el.Snake {
//	x, y := randomPair(gs.WORLD_SIZE)
//	subWorld := wo.SubWorlds()[x][y]
//
//	snake := el.Snake{}
//	subWorldCoord := gs.NewCoord(x, y)
//	gridCoord := placeElementRandomLocationWithLock(&subWorld, &snake)
//
//	snake.SubWorldCoord = subWorldCoord
//	snake.GridCoord = gridCoord
//
//	return snake
//}

func nextCoinId() int {
	coinIdInc++
	return coinIdInc
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
func placeCoin(subWorld *gs.SubWorld, coin *el.Coin) {
	placeElementRandomLocationWithLock(subWorld, coin)
}

// creates a Rock
func placeRock(subWorld *gs.SubWorld) {
	rock := elementFactory.CreateNew(el.ROCK)
	elementFactory.Save(rock)
	placeElementRandomLocationWithLock(subWorld, rock)
}

func placeElementRandomLocationWithLock(subWorld *gs.SubWorld, dbo rc.Dbo) gs.Coord {
	x, y := randomPair(gs.GRID_SIZE)
	coord := gs.NewCoord(x, y)

	if isEmpty(subWorld.Coord(), coord) {
		cell := subWorld.Grid()[x][y]

		cell.Mux().Lock()
		storeElement(subWorld, coord, dbo)
		cell.Mux().Unlock()
	} else {
		coord = placeElementRandomLocationWithLock(subWorld, dbo)
	}

	return coord
}

func movePlayer(wo *gs.World, player *el.Player, vector gs.Vector) {
	player.SubWorldCoord, player.GridCoord = moveCharacter(wo, player.SubWorldCoord, player.GridCoord, vector, player)
}

//func moveSnake(gs *gs.World, snake *el.Snake, vector gs.Vector) {
//	snake.SubWorldCoord, snake.GridCoord = moveCharacter(gs, snake.SubWorldCoord, snake.GridCoord, vector, snake)
//}

func moveCharacter(wo *gs.World, subWorldCoord gs.Coord, coord gs.Coord, vector gs.Vector, element rc.Dbo) (gs.Coord, gs.Coord) {
	subWorld := wo.SubWorlds()[subWorldCoord.X][subWorldCoord.Y]

	nextSubWorldCoord, nextCoord, _ := subWorldMove(subWorldCoord, coord, vector)

	nextSubWorld := wo.SubWorlds()[nextSubWorldCoord.X][nextSubWorldCoord.Y]

	nextElement, _ := nextElement(wo, subWorldCoord, coord, vector)

	//fmt.Printf("game.go: NextElement key: %s, value: %s\n", nextElement.String(), nextElement.Serialize())

	override := false
	switch element.(type) {
	//case *el.Snake:
	//	snake := element.(*el.Snake)
	//	override = snake.Interact(nextElement)
	case *el.Player:
		player := element.(*el.Player)
		override = player.Interact(nextElement)
	default:
		panic("I don't know how to move this type")
	}

	if override {
		prevCell := subWorld.Grid()[coord.X][coord.Y]

		prevCell.Mux().Lock()
		removeCoords(subWorldCoord, coord)
		prevCell.Mux().Unlock()

		nextCell := nextSubWorld.Grid()[nextCoord.X][nextCoord.Y]
		nextCell.Mux().Lock()
		storeElement(&nextSubWorld, nextCoord, element)
		nextCell.Mux().Unlock()
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

func nextElement(wo *gs.World, subWorldCoord gs.Coord, coord gs.Coord, vector gs.Vector) (rc.Dbo, bool) {
	nextSubWorldCoord, nextCoord, moved := subWorldMove(subWorldCoord, coord, vector)
	return elementFromCoords(nextSubWorldCoord, nextCoord), moved
}

func wrap(base int, add int, max int) int {
	sum := base + add

	if sum > 0 {
		return sum % max
	} else {
		return ((sum % max) + max) % max
	}
}

func carry(base int, add int, max int) int {
	sum := base + add

	if sum > 0 {
		return sum / max
	} else {
		return ((sum - max + 1) / max)
	}
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

func isOutOfBound(x int, y int, bound int) bool {
	return x < 0 || y < 0 || x >= bound || y >= bound
}

func printWorld(wo *gs.World, player *el.Player) {
	v := gs.NewVector(-5, -5)
	visionDistance := 11

	for i := 0; i < visionDistance; i++ {
		for j := 0; j < visionDistance; j++ {
			element, valid := nextElement(wo, player.SubWorldCoord, player.GridCoord, v)
			if valid {
				fmt.Printf("%v ", element.String())
			}
			v.X += 1
		}
		v.X = -5
		fmt.Println()
		v.Y += 1
	}
}

//func printDebugWorld(wo *gs.World, player *el.Player) {
//	v := gs.NewCoord(-5, -5)
//	visionDistance := 11
//
//	for i := 0; i < visionDistance; i++ {
//		for j := 0; j < visionDistance; j++ {
//			val, _ := redisClient.Get(coordKey(player.SubWorldCoord, v)).Result()
//			fmt.Printf(val)
//			v.X += 1
//		}
//		v.X = -5
//		fmt.Println()
//		v.Y += 1
//	}
//}

func render(wo *gs.World, player *el.Player) {
	for {
		clearScreen()
		printWorld(wo, player)
		fmt.Println()
		printStat(player)

		if debug {
			//printDebugWorld(wo, player)
			//
			//keys, _ := redisClient.Keys("player:*").Result()
			//for _, key := range keys {
			//	val, _ := redisClient.Get(key).Result()
			//	fmt.Printf("[%v--%v]", key, val)
			//}
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func resetTerminal() {
	cmd := exec.Command("cmd", "/c", "reset")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func clearScreen() {
	//cmd := exec.Command("cmd", "/c", "cls || clear")
	//cmd.Stdout = os.Stdout
	//cmd.Run()

	print("\033[H\033[2J")
}

func abs(n int) int {
	if n < 0 {
		return -n
	} else {
		return n
	}
}

func convertToOneMove(n int) int {
	if n < 0 {
		return -1
	} else {
		return 1
	}
}

//func findPlayer(subWorld *gs.SubWorld) gs.Coord {
//	x := -1
//	y := -1
//
//	for i := 0; i < gs.GRID_SIZE; i++ {
//		for j := 0; j < gs.GRID_SIZE; j++ {
//			//_, isPlayer := subWorld.grid[i][j].element.(*Player)
//			if isPlayer {
//				x = i
//				y = j
//				break
//			}
//		}
//	}
//
//	return gs.Coord{X: x, Y: y}
//}
//
//func spawnSnake(gsWorld *gs.World) {
//	snake := initializeSnake(gsWorld)
//
//	for {
//		if snake.GridCoord.X == -1 && snake.GridCoord.Y == -1 {
//			return
//		}
//
//		playerLocation := findPlayer(&gsWorld.SubWorlds()[snake.SubWorldCoord.X][snake.SubWorldCoord.Y])
//
//		moveVector := gs.Vector{X: 0, Y: 0}
//
//		if isFound(playerLocation) {
//			diffX := playerLocation.X - snake.GridCoord.X
//			diffY := playerLocation.Y - snake.GridCoord.Y
//
//			if abs(diffX) > abs(diffY) {
//				moveVector.X = convertToOneMove(diffX)
//			} else {
//				moveVector.Y = convertToOneMove(diffY)
//			}
//		} else {
//			moveVector = randomVector()
//		}
//
//		moveSnake(gsWorld, &snake, moveVector)
//		time.Sleep(250 * time.Millisecond)
//	}
//}

func isFound(coord gs.Coord) bool {
	return coord.X > 0 && coord.Y > 0
}

func randomPair(n int) (int, int) {
	return rand.Intn(n), rand.Intn(n)
}

func randomVector() gs.Vector {
	x := rand.Intn(3) - 1
	y := rand.Intn(3) - 1

	return gs.NewVector(x, y)
}

func randomSubWorldCoord() gs.Coord {
	x, y := randomPair(gs.WORLD_SIZE)

	return gs.NewCoord(x, y)
}

func randomSubWorld(wo *gs.World) gs.SubWorld {
	coord := randomSubWorldCoord()

	return wo.SubWorlds()[coord.X][coord.Y]
}

func spawnCoinsInWorld(wo *gs.World) {
	sleepTime := 10000 * time.Millisecond

	for {
		randomSubWorld := randomSubWorld(wo)
		spawnCoinInSubWorld(&randomSubWorld)
		time.Sleep(sleepTime)
		sleepTime += sleepTime
	}
}

func spawnCoinInSubWorld(subWorld *gs.SubWorld) {
	coin := buildAndStoreCoin()
	placeCoin(subWorld, coin)
}

func printStat(player *el.Player) {
	fmt.Printf("Coin: %d", player.CoinCount)
	fmt.Println()
	fmt.Printf("HP: %d", player.Hp)
	fmt.Println()
}

func checkAlive(player *el.Player) {
	if !player.Alive {
		fmt.Println("You Died")
		os.Exit(0)
	}
}

func startTerminalClient(wo *gs.World) {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	player := initializePlayer(wo)

	go render(wo, player)

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyCtrlQ {
				elementFactory.Save(player)
				//TODO - save globalId
				termbox.Close()
				clearScreen()
				os.Exit(3)
			}

			moveVector := gs.NewVector(0, 0)
			if ev.Ch == 119 { // w
				moveVector.Y = -1
				movePlayer(wo, player, moveVector)
			} else if ev.Ch == 97 { // a
				moveVector.X = -1
				movePlayer(wo, player, moveVector)
			} else if ev.Ch == 115 { // s
				moveVector.Y = 1
				movePlayer(wo, player, moveVector)
			} else if ev.Ch == 100 { // d
				moveVector.X = 1
				movePlayer(wo, player, moveVector)
				//} else if ev.Ch == 0 { // space
				//	player.BuildWall(gs)
			}

			checkAlive(player)

			termbox.Flush()
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}

func initializeWorld() gs.World {
	subWorlds := [gs.WORLD_SIZE][gs.WORLD_SIZE]gs.SubWorld{}

	for i := 0; i < gs.WORLD_SIZE; i++ {
		for j := 0; j < gs.WORLD_SIZE; j++ {
			subWorlds[i][j] = initializeSubWorld(i, j)
		}
	}

	return gs.NewWorld(subWorlds)
}

func initializeSubWorld(i int, j int) gs.SubWorld {
	subWorld := gs.NewSubWorld(gs.NewCoord(i, j))

	if firstBoot {
		for i := 0; i < 10; i++ {
			placeRock(&subWorld)
		}

		for i := 0; i < 2; i++ {
			coin := buildAndStoreCoin()
			placeCoin(&subWorld, coin)
		}
	}

	return subWorld
}

//func spawnSnakes(gsWorld *World) {
//	for {
//		go spawnSnake(gsWorld)
//		time.Sleep(3000 * time.Millisecond)
//	}
//}

func initializeElementFactory() {
	elementFactory = el.Factory(firstBoot)
	elementFactory.Init()
}

func runWorldElements(gs *gs.World) {
	//go spawnSnakes(gs)
	go spawnCoinsInWorld(gs)
}

func main() {
	rand.Seed(12345)

	firstBoot = false
	debug = true

	initializeElementFactory()

	wo := initializeWorld()

	runWorldElements(&wo)

	startTerminalClient(&wo)

	// TESTS
	//fmt.Println(carry(0, 3, 3) == 1)
	//fmt.Println(carry(0, 4, 3) == 1)
	//fmt.Println(carry(0, 5, 3) == 1)
	//fmt.Println(carry(0, 6, 3) == 2)
	//fmt.Println(carry(0, 1, 3) == 0)
	//fmt.Println(carry(0, 2, 3) == 0)
	//fmt.Println(carry(0, 3, 3) == 1)
	//fmt.Println(carry(0, -1, 3) == -1)
	//fmt.Println(carry(0, -2, 3) == -1)
	//fmt.Println(carry(0, -3, 3) == -1)
	//fmt.Println(carry(0, -4, 3) == -2)
	//fmt.Println(carry(0, -7, 3) == -3)
	//fmt.Println(wrap(0, -1, 4) == 3)
	//fmt.Println(wrap(0, -2, 4) == 2)
	//fmt.Println(wrap(0, -3, 4) == 1)
	//fmt.Println(wrap(0, -4, 4) == 0)
}
*/
