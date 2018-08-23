package main

import "github.com/nsf/termbox-go"
import "github.com/go-redis/redis"
import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)
import (
	"./el"
	"./gs"
)

const MAX_COIN_AMOUNT = 10

func initializePlayerFromValues(elementId string, values string) el.Player {
	keyValues := strings.Split(values, "-")

	idString := strings.Split(elementId, "player:")[1]

	coinCountString := strings.Split(keyValues[0], "coinCount:")[1]
	aliveString := strings.Split(keyValues[1], "alive:")[1]
	hpString := strings.Split(keyValues[2], "hp:")[1]

	subWorldCoordString := strings.Split(keyValues[3], "subWorldCoord:")[1]
	subWorldCoordStringX := strings.Split(subWorldCoordString, ",")[0]
	subWorldCoordX, _ := strconv.Atoi(subWorldCoordStringX)
	subWorldCoordStringY := strings.Split(subWorldCoordString, ",")[1]
	subWorldCoordY, _ := strconv.Atoi(subWorldCoordStringY)

	gridCoordString := strings.Split(keyValues[4], "gridCoord:")[1]
	gridCoordStringX := strings.Split(gridCoordString, ",")[0]
	gridCoordX, _ := strconv.Atoi(gridCoordStringX)
	gridCoordStringY := strings.Split(gridCoordString, ",")[1]
	gridCoordY, _ := strconv.Atoi(gridCoordStringY)

	id, _ := strconv.Atoi(idString)
	coinCount, _ := strconv.Atoi(coinCountString)
	hp, _ := strconv.Atoi(hpString)
	alive := aliveString == "true"

	player := el.NewPlayer(id, coinCount, alive, hp)

	player.SubWorldCoord = gs.NewCoord(subWorldCoordX, subWorldCoordY)
	player.GridCoord = gs.NewCoord(gridCoordX, gridCoordY)

	return player
}

func initializeCoinFromValues(elementId string, values string) *el.Coin {
	amountString := strings.Split(values, "amount:")[1]
	amount, _ := strconv.Atoi(amountString)

	idString := strings.Split(elementId, "player:")[1]
	id, _ := strconv.Atoi(idString)

	return el.NewCoin(amount, id)
}

func initializeRockFromValues(_ string) el.Rock {
	return el.NewRock()
}

// GLOBALS
var redisClient *redis.Client
var firstBoot bool
var debug bool
var coinIdInc int

func redisSet(key string, value string) {
	err := redisClient.Set(key, value, 0).Err()
	if err != nil {
		panic(err)
	}
}
func storeObject(id string, object string) {
	redisSet(id, object)
}

func storePlayer(player *el.Player) {
	storeObject(player.Id(), player.Val())
}
func storeCoin(coin *el.Coin) {
	storeObject(coin.Id(), coin.Val())
}

func coordKey(subWorldCoord gs.Coord, coord gs.Coord) string {
	return fmt.Sprintf("%v:%v", subWorldCoord.Key(), coord.Key())
}
func subWorldCoordKey(subWorld *gs.SubWorld, coord gs.Coord) string {
	return coordKey(subWorld.Coord(), coord)
}

func storeCoord(subWorld *gs.SubWorld, coord gs.Coord, id string) {
	redisSet(subWorldCoordKey(subWorld, coord), id)
}

func storeCoinCoord(subWorld *gs.SubWorld, coord gs.Coord, coin *el.Coin) {
	storeCoord(subWorld, coord, coin.Id())
}
func storeRockCoord(subWorld *gs.SubWorld, coord gs.Coord, rock *el.Rock) {
	storeCoord(subWorld, coord, rock.Id())
}
func storePlayerCoord(subWorld *gs.SubWorld, coord gs.Coord, player *el.Player) {
	storeCoord(subWorld, coord, player.Id())
}

func emptykey(key string) {
	err := redisClient.Set(key, nil, 0).Err()
	if err != nil {
		panic(err)
	}
}

func setEmptyValue(subWorldCoord gs.Coord, coord gs.Coord) {
	emptykey(coordKey(subWorldCoord, coord))
}

func setEmptyObject(id string) {
	emptykey(id)
}

func isEmpty(coord gs.Coord) bool {
	// Look into this, not sure this is right
	// can we check that value == nil (?)

	val, _ := redisClient.Get(coord.Key()).Result()
	if val == "" {
		return true
	} else {
		return false
	}
}

func storeElement(subWorld *gs.SubWorld, coord gs.Coord, element interface{}) {
	switch v := element.(type) {
	case *el.Coin:
		storeCoinCoord(subWorld, coord, v)
	case *el.Rock:
		storeRockCoord(subWorld, coord, v)
	case *el.Player:
		storePlayerCoord(subWorld, coord, v)
		storePlayer(v)
		//default:
	}
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
func initializePlayer(world *gs.World) *el.Player {
	x, y := randomPair(gs.WORLD_SIZE)
	subWorld := world.SubWorlds()[x][y]

	var player el.Player
	if firstBoot {
		player = el.NewPlayer(1, 0, true, 10)

		subWorldCoord := gs.NewCoord(x, y)
		gridCoord := placeElementRandomLocationWithLock(&subWorld, player)

		player.SubWorldCoord = subWorldCoord
		player.GridCoord = gridCoord

		storePlayer(&player)
	} else {
		playerValues, err := redisClient.Get(playerKey(1)).Result()
		if err != nil {
			panic(err)
		}

		player = initializePlayerFromValues("player:1", playerValues)
	}

	return &player
}

// creates a snakeCell
// returns a pair of Coord of World, SubWorld
func initializeSnake(world *gs.World) el.Snake {
	x, y := randomPair(gs.WORLD_SIZE)
	subWorld := world.SubWorlds()[x][y]

	snake := el.Snake{}
	subWorldCoord := gs.NewCoord(x, y)
	gridCoord := placeElementRandomLocationWithLock(&subWorld, &snake)

	snake.SubWorldCoord = subWorldCoord
	snake.GridCoord = gridCoord

	return snake
}

func nextCoinId() int {
	coinIdInc++
	return coinIdInc
}

func initializeCoin() *el.Coin {
	return el.NewCoin(rand.Intn(MAX_COIN_AMOUNT)+1, nextCoinId())
}

func buildAndStoreCoin() *el.Coin {
	coin := initializeCoin()
	storeCoin(coin)
	return coin
}

// creates a Coin
// returns a pair of Coord of SubWorld
func placeCoin(subWorld *gs.SubWorld, coin *el.Coin) {
	placeElementRandomLocationWithLock(subWorld, coin)
}

// creates a Rock
// returns a pair of Coord of SubWorld
func placeRock(subWorld *gs.SubWorld) {
	rock := el.Rock{}
	placeElementRandomLocationWithLock(subWorld, &rock)
}

func placeElementRandomLocationWithLock(subWorld *gs.SubWorld, element interface{}) gs.Coord {
	x, y := randomPair(gs.GRID_SIZE)
	coord := gs.NewCoord(x, y)

	if isEmpty(coord) {
		cell := subWorld.Grid()[x][y]

		cell.Mux().Lock()
		storeElement(subWorld, coord, element)
		cell.Mux().Unlock()
	} else {
		coord = placeElementRandomLocationWithLock(subWorld, element)
	}

	return coord
}

func movePlayer(world *gs.World, player *el.Player, vector gs.Vector) {
	player.SubWorldCoord, player.GridCoord = moveCharacter(world, player.SubWorldCoord, player.GridCoord, vector, player)
}

func moveSnake(gs *gs.World, snake *el.Snake, vector gs.Vector) {
	snake.SubWorldCoord, snake.GridCoord = moveCharacter(gs, snake.SubWorldCoord, snake.GridCoord, vector, snake)
}

func moveCharacter(world *gs.World, subWorldCoord gs.Coord, coord gs.Coord, vector gs.Vector, element interface{}) (gs.Coord, gs.Coord) {
	subWorld := world.SubWorlds()[subWorldCoord.X][subWorldCoord.Y]

	nextSubWorldCoord, nextCoord, _ := subWorldMove(subWorldCoord, coord, vector)

	nextSubWorld := world.SubWorlds()[nextSubWorldCoord.X][nextSubWorldCoord.Y]

	nextElement, _ := nextElement(world, subWorldCoord, coord, vector)

	override := false
	switch element.(type) {
	case *el.Snake:
		snake := element.(*el.Snake)
		override = snake.Interact(nextElement)
	case *el.Player:
		player := element.(*el.Player)
		override = player.Interact(nextElement)
	default:
		panic("I don't know how to move this type")
	}

	if override {
		prevCell := subWorld.Grid()[coord.X][coord.Y]

		prevCell.Mux().Lock()
		setEmptyValue(subWorldCoord, coord)
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

func elementFromElementId(elementId string) interface{} {
	elementValues, _ := redisClient.Get(elementId).Result()

	var element interface{}
	switch strings.Split(elementId, ":")[0] {
	case "coin":
		element = initializeCoinFromValues(elementId, elementValues)
	case "rock":
		element = initializeRockFromValues(elementValues)
	case "player":
		element = initializePlayerFromValues(elementId, elementValues)
		//case "snake":
		//	element = initializeSnakeFromValues(elementValues)
	}

	return element
}

func elementFromCoords(nextSubWorldCoord gs.Coord, nextCoord gs.Coord) interface{} {
	elementId, _ := redisClient.Get(coordKey(nextSubWorldCoord, nextCoord)).Result()

	if elementId == "" {
		return el.Empty{}
	} else {
		return elementFromElementId(elementId)
	}
}

func nextElement(world *gs.World, subWorldCoord gs.Coord, coord gs.Coord, vector gs.Vector) (interface{}, bool) {
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

func printWorld(world *gs.World, player *el.Player) {
	v := gs.NewVector(-5, -5)
	visionDistance := 11

	for i := 0; i < visionDistance; i++ {
		for j := 0; j < visionDistance; j++ {
			element, valid := nextElement(world, player.SubWorldCoord, player.GridCoord, v)
			if valid {
				fmt.Printf("%v ", element)
			}
			v.X += 1
		}
		v.X = -5
		fmt.Println()
		v.Y += 1
	}
}

func printDebugWorld(world *gs.World, player *el.Player) {
	v := gs.NewCoord(-5, -5)
	visionDistance := 11

	for i := 0; i < visionDistance; i++ {
		for j := 0; j < visionDistance; j++ {
			val, _ := redisClient.Get(coordKey(player.SubWorldCoord, v)).Result()
			fmt.Printf(val)
			v.X += 1
		}
		v.X = -5
		fmt.Println()
		v.Y += 1
	}
}

func render(world *gs.World, player *el.Player) {
	for {
		clearScreen()
		printWorld(world, player)
		fmt.Println()

		if debug {
			printDebugWorld(world, player)
			printStat(player)

			keys, _ := redisClient.Keys("player:*").Result()
			for _, key := range keys {
				val, _ := redisClient.Get(key).Result()
				fmt.Printf("[%v--%v]", key, val)
			}
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

func randomSubWorld(world *gs.World) gs.SubWorld {
	coord := randomSubWorldCoord()

	return world.SubWorlds()[coord.X][coord.Y]
}

func spawnCoinsInWorld(world *gs.World) {
	sleepTime := 10000 * time.Millisecond

	for {
		randomSubWorld := randomSubWorld(world)
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
	fmt.Printf("Coin: %d", player.CoinCount())
	fmt.Println()
	fmt.Printf("HP: %d", player.Hp())
	fmt.Println()
}

func checkAlive(player *el.Player) {
	if !player.Alive() {
		fmt.Println("You Died")
		os.Exit(0)
	}
}

func startTerminalClient(world *gs.World) {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	player := initializePlayer(world)

	go render(world, player)

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyCtrlQ {
				storePlayer(player)
				termbox.Close()
				clearScreen()
				os.Exit(3)
			}

			moveVector := gs.NewVector(0, 0)
			if ev.Ch == 119 { // w
				moveVector.Y = -1
				movePlayer(world, player, moveVector)
			} else if ev.Ch == 97 { // a
				moveVector.X = -1
				movePlayer(world, player, moveVector)
			} else if ev.Ch == 115 { // s
				moveVector.Y = 1
				movePlayer(world, player, moveVector)
			} else if ev.Ch == 100 { // d
				moveVector.X = 1
				movePlayer(world, player, moveVector)
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

func runWorldElements(gs *gs.World) {
	//go spawnSnakes(gs)
	go spawnCoinsInWorld(gs)
}

func main() {
	rand.Seed(12345)

	firstBoot = true
	debug = true

	initializeRedisClient()

	world := initializeWorld()

	runWorldElements(&world)

	startTerminalClient(&world)

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
