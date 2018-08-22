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
	"sync"
	"time"
)
import "./gs"

const MAX_COIN_AMOUNT = 10

type Element struct {
	subWorldCoord gs.Coord
	gridCoord     gs.Coord
}

// PLAYER
type Player struct {
	Element
	mux       sync.Mutex
	coinCount int
	alive     bool
	hp        int
	id        int
}

func (p Player) String() string {
	return fmt.Sprintf("%v", p.id)
}
func (player *Player) Interact(element interface{}) bool {
	switch v := element.(type) {
	case Coin:
		player.IncCoinCount(v.amount)
		v.Destroy()
		return true
	case Empty:
		return true
	default:
		return false
	}

	return false
}
func (player *Player) Kill() {
	player.mux.Lock()
	player.alive = false
	player.mux.Unlock()
}
func (player *Player) IncCoinCount(amount int) {
	player.mux.Lock()
	player.coinCount += amount
	storePlayer(player)
	player.mux.Unlock()
}
func (player *Player) decreaseHp(damage int) {
	player.mux.Lock()
	player.hp -= damage
	player.mux.Unlock()

	if player.hp < 0 {
		player.Kill()
	}
}
func (player *Player) Id() string {
	return fmt.Sprintf("player:%v", player.id)
}
func (player *Player) Val() string {
	// Bug fix, use dashes because cords use commas. FIXME: use commas for all attr delimiters
	return fmt.Sprintf("coinCount:%v-alive:%v-hp:%v-subWorldCoord:%v-gridCoord:%v",
		player.coinCount,
		player.alive,
		player.hp,
		player.subWorldCoord.Key(),
		player.gridCoord.Key(),
	)
}

func initializePlayerFromValues(elementId string, values string) Player {
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

	player := Player{id: id, coinCount: coinCount, alive: alive, hp: hp}

	player.subWorldCoord = gs.NewCoord(subWorldCoordX, subWorldCoordY)
	player.gridCoord = gs.NewCoord(gridCoordX, gridCoordY)

	return player
}

// SNAKE
type Snake struct {
	Element
}

func (s *Snake) String() string {
	return "S"
}
func (s *Snake) Attack(player *Player) {
	player.decreaseHp(1)
}
func (snake *Snake) Interact(element interface{}) bool {
	switch element.(type) {
	case *Player:
		player := element.(*Player)
		snake.Attack(player)
		return true
	case *Empty:
		return true
	default:
		return false
	}

	return false
}

// COIN
type Coin struct {
	Element
	amount int
	id     int
}

func (c Coin) String() string {
	return "C"
}
func (coin *Coin) Id() string {
	return fmt.Sprintf("coin:%v", coin.id)
}
func (coin *Coin) Val() string {
	return fmt.Sprintf("amount:%v", coin.amount)
}
func (coin *Coin) Destroy() {
	setEmptyObject(coin.Id())
}

func initializeCoinFromValues(values string) Coin {
	amountString := strings.Split(values, "amount:")[1]
	amount, _ := strconv.Atoi(amountString)
	return Coin{amount: amount}
}

// ROCK
type Rock struct {
	Element
}

func (rock *Rock) Id() string {
	return fmt.Sprintf("rock")
}
func (r Rock) String() string {
	return "R"
}

func initializeRockFromValues(_ string) Rock {
	return Rock{}
}

// EMPTY
type Empty struct {
	Element
}

func (e Empty) String() string {
	return " "
}

// BUILDING
type Building struct {
	Element
	code string
}

func (building Building) String() string {
	return building.code
}

// GLOBALS
var redisClient *redis.Client
var firstBoot bool
var debug bool
var coinIdInc int

func initializeRedisClient() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func redisSet(key string, value string) {
	err := redisClient.Set(key, value, 0).Err()
	if err != nil {
		panic(err)
	}
}
func storeObject(id string, object string) {
	redisSet(id, object)
}

func storePlayer(player *Player) {
	storeObject(player.Id(), player.Val())
}
func storeCoin(coin *Coin) {
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

func storeCoinCoord(subWorld *gs.SubWorld, coord gs.Coord, coin *Coin) {
	storeCoord(subWorld, coord, coin.Id())
}
func storeRockCoord(subWorld *gs.SubWorld, coord gs.Coord, rock *Rock) {
	storeCoord(subWorld, coord, rock.Id())
}
func storePlayerCoord(subWorld *gs.SubWorld, coord gs.Coord, player *Player) {
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
	case *Coin:
		storeCoinCoord(subWorld, coord, v)
	case *Rock:
		storeRockCoord(subWorld, coord, v)
	case *Player:
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
//	upCell := nextCell(gs, player.subWorldCoord, player.gridCoord, up)
//	upElement := upCell.element
//
//	upRightCell := nextCell(gs, player.subWorldCoord, player.gridCoord, upRight)
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
func initializePlayer(world *gs.World) *Player {
	x, y := randomPair(gs.WORLD_SIZE)
	subWorld := world.SubWorlds()[x][y]

	var player Player
	if firstBoot {
		player = Player{id: 1, alive: true, hp: 10}

		subWorldCoord := gs.NewCoord(x, y)
		gridCoord := placeElementRandomLocationWithLock(&subWorld, player)

		player.subWorldCoord = subWorldCoord
		player.gridCoord = gridCoord

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
func initializeSnake(world *gs.World) Snake {
	x, y := randomPair(gs.WORLD_SIZE)
	subWorld := world.SubWorlds()[x][y]

	snake := Snake{}
	subWorldCoord := gs.NewCoord(x, y)
	gridCoord := placeElementRandomLocationWithLock(&subWorld, &snake)

	snake.subWorldCoord = subWorldCoord
	snake.gridCoord = gridCoord

	return snake
}

func nextCoinId() int {
	coinIdInc++
	return coinIdInc
}

func initializeCoin() *Coin {
	return &Coin{amount: rand.Intn(MAX_COIN_AMOUNT) + 1, id: nextCoinId()}
}

func buildAndStoreCoin() *Coin {
	coin := initializeCoin()
	storeCoin(coin)
	return coin
}

// creates a Coin
// returns a pair of Coord of SubWorld
func placeCoin(subWorld *gs.SubWorld, coin *Coin) {
	placeElementRandomLocationWithLock(subWorld, coin)
}

// creates a Rock
// returns a pair of Coord of SubWorld
func placeRock(subWorld *gs.SubWorld) {
	rock := Rock{}
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

func movePlayer(world *gs.World, player *Player, vector gs.Vector) {
	player.subWorldCoord, player.gridCoord = moveCharacter(world, player.subWorldCoord, player.gridCoord, vector, player)
}

//func moveSnake(gs *World, snake *Snake, vector Vector) {
//	snake.subWorldCoord, snake.gridCoord = moveCharacter(gs, snake.subWorldCoord, snake.gridCoord, vector, snake)
//}
func moveCharacter(world *gs.World, subWorldCoord gs.Coord, coord gs.Coord, vector gs.Vector, element interface{}) (gs.Coord, gs.Coord) {
	subWorld := world.SubWorlds()[subWorldCoord.X][subWorldCoord.Y]

	nextSubWorldCoord, nextCoord, _ := subWorldMove(subWorldCoord, coord, vector)

	nextSubWorld := world.SubWorlds()[nextSubWorldCoord.X][nextSubWorldCoord.Y]

	nextElement, _ := nextElement(world, subWorldCoord, coord, vector)

	override := false
	switch element.(type) {
	case *Snake:
		snake := element.(*Snake)
		override = snake.Interact(nextElement)
	case *Player:
		player := element.(*Player)
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
		element = initializeCoinFromValues(elementValues)
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
		return Empty{}
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

func printWorld(world *gs.World, player *Player) {
	v := gs.NewVector(-5, -5)
	visionDistance := 11

	for i := 0; i < visionDistance; i++ {
		for j := 0; j < visionDistance; j++ {
			element, valid := nextElement(world, player.subWorldCoord, player.gridCoord, v)
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

func printDebugWorld(world *gs.World, player *Player) {
	v := gs.NewCoord(-5, -5)
	visionDistance := 11

	for i := 0; i < visionDistance; i++ {
		for j := 0; j < visionDistance; j++ {
			val, _ := redisClient.Get(coordKey(player.subWorldCoord, v)).Result()
			fmt.Printf(val)
			v.X += 1
		}
		v.X = -5
		fmt.Println()
		v.Y += 1
	}
}

func render(world *gs.World, player *Player) {
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

//func findPlayer(subWorld *gs.SubWorld) Coord {
//	x := -1
//	y := -1
//
//	for i := 0; i < GRID_SIZE; i++ {
//		for j := 0; j < GRID_SIZE; j++ {
//			_, isPlayer := subWorld.grid[i][j].element.(*Player)
//			if isPlayer {
//				x = i
//				y = j
//				break
//			}
//		}
//	}
//
//	return Coord{x: x, y: y}
//}

//func spawnSnake(gs *World) {
//	snake := initializeSnake(gs)
//
//	for {
//		if snake.gridCoord.x == -1 && snake.gridCoord.y == -1 {
//			return
//		}
//
//		playerLocation := findPlayer(&gs.subWorlds[snake.subWorldCoord.x][snake.subWorldCoord.y])
//
//		moveVector := Vector{x: 0, y: 0}
//
//		if isFound(playerLocation) {
//			diffX := playerLocation.x - snake.gridCoord.x
//			diffY := playerLocation.y - snake.gridCoord.y
//
//			if abs(diffX) > abs(diffY) {
//				moveVector.x = convertToOneMove(diffX)
//			} else {
//				moveVector.y = convertToOneMove(diffY)
//			}
//		} else {
//			moveVector = randomVector()
//		}
//
//		moveSnake(gs, &snake, moveVector)
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

func printStat(player *Player) {
	fmt.Printf("Coin: %d", player.coinCount)
	fmt.Println()
	fmt.Printf("HP: %d", player.hp)
	fmt.Println()
}

func checkAlive(player *Player) {
	if !player.alive {
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

		// snake := buildSnake()
		// placeSnake(&subWorld, snake)
	}

	return subWorld
}

//func spawnSnakes(gs *World) {
//	for {
//		go spawnSnake(gs)
//		time.Sleep(3000 * time.Millisecond)
//	}
//}

func runWorldElements(world *gs.World) {
	//go spawnSnakes(gs)
	go spawnCoinsInWorld(world)
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
