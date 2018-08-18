package main

import "github.com/nsf/termbox-go"
import "github.com/go-redis/redis"
import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sync"
	"time"
)

const GRID_SIZE = 8
const WORLD_SIZE = 4
const MAX_COIN_AMOUNT = 10

type World struct {
	subWorlds [WORLD_SIZE][WORLD_SIZE]SubWorld
}
type Grid [GRID_SIZE][GRID_SIZE]Cell
type SubWorld struct {
	grid Grid
}
type Coord struct {
	x int
	y int
}
type Vector struct {
	x int
	y int
}
type Cell struct {
	mux     sync.Mutex
	element interface{}
}

func (coord *Coord) Key() string {
	return fmt.Sprintf("%v,%v", coord.x, coord.y)
}

type Element struct {
	subWorldCoord Coord
	gridCoord     Coord
}

// PLAYER
type Player struct {
	Element
	mux       sync.Mutex
	coinCount int
	alive     bool
	hp        int
}

func (p *Player) String() string {
	return "P"
}
func (player *Player) Interact(element interface{}) bool {
	switch v := element.(type) {
	case *Coin:
		player.IncCoinCount(v.amount)
		return true
	case *Empty:
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

func (c *Coin) String() string {
	return "C"
}

func (coin *Coin) Id() string {
	return fmt.Sprintf("coin:%v", coin.id)
}

func (coin *Coin) Val() string {
	return fmt.Sprintf("amount:%v", coin.amount)
}

// ROCK
type Rock struct {
	Element
}

func (rock *Rock) Id() string {
	return fmt.Sprintf("rock")
}

func (r *Rock) String() string {
	return "R"
}

// EMPTY
type Empty struct {
	Element
}

func (e *Empty) String() string {
	return " "
}

type Building struct {
	Element
	code string
}

func (building *Building) String() string {
	return building.code
}

// GLOBALS
var redisClient *redis.Client
var firstBoot bool
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

func storeCoin(coin *Coin) {
	storeObject(coin.Id(), coin.Val())
}

func storeCoord(coord Coord, id string) {
	redisSet(coord.Key(), id)
}

func storeCoinCoord(coord Coord, coin *Coin) {
	storeCoord(coord, coin.Id())
}

func storeRockCoord(coord Coord, rock *Rock) {
	storeCoord(coord, rock.Id())
}

func isEmpty(coord Coord) bool {
	_, err := redisClient.Get(coord.Key()).Result()
	if err == nil {
		return false
	} else {
		return true
	}
}

func storeElement(coord Coord, element interface{}) {
	switch v := element.(type) {
	case *Coin:
		storeCoinCoord(coord, v)
	case *Rock:
		storeRockCoord(coord, v)
		//default:
	}
}

func (player *Player) BuildWall(world *World) bool {
	wallCost := 5

	if player.coinCount < wallCost {
		return false
	}

	up := Vector{x: 0, y: -1}
	upRight := Vector{x: 1, y: -1}

	upCell := nextCell(world, player.subWorldCoord, player.gridCoord, up)
	upElement := upCell.element

	upRightCell := nextCell(world, player.subWorldCoord, player.gridCoord, upRight)
	upRightElement := upRightCell.element

	if isEmpty(upElement) && isEmpty(upRightElement) {
		player.mux.Lock()
		player.coinCount -= wallCost
		player.mux.Unlock()

		upCell.mux.Lock()
		upCell.element = &Building{code: "<"}
		upCell.mux.Unlock()

		upRightCell.mux.Lock()
		upRightCell.element = &Building{code: ">"}
		upRightCell.mux.Unlock()

		return true
	}

	return false
}

// creates a player
// returns a pair of Coord of World, SubWorld
func initializePlayer(world *World) Player {
	x, y := randomPair(WORLD_SIZE)
	subWorld := world.subWorlds[x][y]

	player := Player{alive: true, hp: 10}
	subWorldCoord := Coord{x: x, y: y}
	gridCoord := placeElementRandomLocationWithLock(&subWorld, &player)

	player.subWorldCoord = subWorldCoord
	player.gridCoord = gridCoord

	return player
}

// creates a snakeCell
// returns a pair of Coord of World, SubWorld
func initializeSnake(world *World) Snake {
	x, y := randomPair(WORLD_SIZE)
	subWorld := world.subWorlds[x][y]

	snake := Snake{}
	subWorldCoord := Coord{x: x, y: y}
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
func placeCoin(subWorld *SubWorld, coin *Coin) {
	placeElementRandomLocationWithLock(subWorld, coin)
}

// creates a Rock
// returns a pair of Coord of SubWorld
func placeRock(subWorld *SubWorld) {
	rock := Rock{}
	placeElementRandomLocationWithLock(subWorld, &rock)
}

func placeElementRandomLocationWithLock(subWorld *SubWorld, element interface{}) (Coord, error) {
	x, y := randomPair(GRID_SIZE)
	coord := Coord{x: x, y: y}

	if isEmpty(coord) {
		cell := &subWorld.grid[x][y]

		cell.mux.Lock()
		storeElement(coord, element)
		cell.mux.Unlock()

		return coord, nil
	} else {
		return coord, errors.New("No Space available")
	}
}

func movePlayer(world *World, player *Player, vector Vector) {
	player.subWorldCoord, player.gridCoord = moveCharacter(world, player.subWorldCoord, player.gridCoord, vector, player)
}

func moveSnake(world *World, snake *Snake, vector Vector) {
	snake.subWorldCoord, snake.gridCoord = moveCharacter(world, snake.subWorldCoord, snake.gridCoord, vector, snake)
}

func moveCharacter(world *World, subWorldCoord Coord, coord Coord, vector Vector, element interface{}) (Coord, Coord) {
	subWorld := &world.subWorlds[subWorldCoord.x][subWorldCoord.y]

	nextSubWorldCoord, nextCoord := subWorldMove(subWorldCoord, coord, vector)

	nextSubWorld := &world.subWorlds[nextSubWorldCoord.x][nextSubWorldCoord.y]

	nextElement := nextElement(world, subWorldCoord, coord, vector)

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
		prevCell := &subWorld.grid[coord.x][coord.y]

		prevCell.mux.Lock()
		// should we put coords in this empty cell?
		prevCell.element = &Empty{}
		prevCell.mux.Unlock()

		nextCell := &nextSubWorld.grid[nextCoord.x][nextCoord.y]
		nextCell.mux.Lock()
		nextCell.element = element
		nextCell.mux.Unlock()
	} else {
		nextSubWorldCoord = subWorldCoord
		nextCoord = coord
	}
	return nextSubWorldCoord, nextCoord
}

func nextCell(world *World, subWorldCoord Coord, coord Coord, vector Vector) *Cell {
	nextSubWorldCoord, nextCoord := subWorldMove(subWorldCoord, coord, vector)
	nextSubWorld := &world.subWorlds[nextSubWorldCoord.x][nextSubWorldCoord.y]
	nextCell := getCell(nextSubWorld, nextCoord)

	return nextCell
}

func nextElement(world *World, subWorldCoord Coord, coord Coord, vector Vector) interface{} {
	return nextCell(world, subWorldCoord, coord, vector).element
}

func getCell(subWorld *SubWorld, coord Coord) *Cell {
	return &subWorld.grid[coord.x][coord.y]
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

func subWorldMove(subWorldCoord Coord, gridCoord Coord, vector Vector) (Coord, Coord) {
	wX := subWorldCoord.x + carry(gridCoord.x, vector.x, GRID_SIZE)
	wY := subWorldCoord.y + carry(gridCoord.y, vector.y, GRID_SIZE)

	gX := wrap(gridCoord.x, vector.x, GRID_SIZE)
	gY := wrap(gridCoord.y, vector.y, GRID_SIZE)

	if isOutOfBound(wX, wY, WORLD_SIZE) {
		return subWorldCoord, gridCoord
	}

	return Coord{x: wX, y: wY}, Coord{x: gX, y: gY}
}

func isOutOfBound(x int, y int, bound int) bool {
	return x < 0 || y < 0 || x >= bound || y >= bound
}

func printWorld(world *World, player *Player) {
	v := Vector{x: -5, y: -5}
	visionDistance := 11

	for i := 0; i < visionDistance; i++ {
		for j := 0; j < visionDistance; j++ {
			fmt.Printf("%v ", nextCell(world, player.subWorldCoord, player.gridCoord, v).element)
			//fmt.Print(v)
			v.x += 1
		}
		v.x = -5
		fmt.Println()
		v.y += 1
	}
}

func render(world *World, player *Player) {
	for {
		clearScreen()
		printWorld(world, player)
		printStat(player)
		time.Sleep(100 * time.Millisecond)
	}
}

func clearScreen() {
	cmd := exec.Command("cmd", "/c", "cls || clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	//print("\033[H\033[2J")
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

func findPlayer(subWorld *SubWorld) Coord {
	x := -1
	y := -1

	for i := 0; i < GRID_SIZE; i++ {
		for j := 0; j < GRID_SIZE; j++ {
			_, isPlayer := subWorld.grid[i][j].element.(*Player)
			if isPlayer {
				x = i
				y = j
				break
			}
		}
	}

	return Coord{x: x, y: y}
}

func spawnSnake(world *World) {
	snake := initializeSnake(world)

	for {
		if snake.gridCoord.x == -1 && snake.gridCoord.y == -1 {
			return
		}

		playerLocation := findPlayer(&world.subWorlds[snake.subWorldCoord.x][snake.subWorldCoord.y])

		moveVector := Vector{x: 0, y: 0}

		if isFound(playerLocation) {
			diffX := playerLocation.x - snake.gridCoord.x
			diffY := playerLocation.y - snake.gridCoord.y

			if abs(diffX) > abs(diffY) {
				moveVector.x = convertToOneMove(diffX)
			} else {
				moveVector.y = convertToOneMove(diffY)
			}
		} else {
			moveVector = randomVector()
		}

		moveSnake(world, &snake, moveVector)
		time.Sleep(250 * time.Millisecond)
	}
}

func isFound(coord Coord) bool {
	return coord.x > 0 && coord.y > 0
}

func randomPair(n int) (int, int) {
	return rand.Intn(n), rand.Intn(n)
}

func randomVector() Vector {
	x := rand.Intn(3) - 1
	y := rand.Intn(3) - 1

	return Vector{x: x, y: y}
}

func randomSubWorldCoord() Coord {
	x, y := randomPair(WORLD_SIZE)

	return Coord{x: x, y: y}
}

func randomSubWorld(world *World) *SubWorld {
	coord := randomSubWorldCoord()

	return &world.subWorlds[coord.x][coord.y]
}

func spawnCoinsInWorld(world *World) {
	sleepTime := 10000 * time.Millisecond

	for {
		randomSubWorld := randomSubWorld(world)
		spawnCoinInSubWorld(randomSubWorld)
		time.Sleep(sleepTime)
		sleepTime += sleepTime
	}
}

func spawnCoinInSubWorld(subWorld *SubWorld) {
	coin := initializeCoin()
	placeCoin(subWorld, coin)
}

func printStat(player *Player) {
	fmt.Printf("Coin: %d", player.coinCount)
	fmt.Println()
	fmt.Printf("HP: %d", player.hp)
}

func checkAlive(player *Player) {
	if !player.alive {
		fmt.Println("You Died")
		os.Exit(0)
	}
}

func startTerminalClient(world *World) {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	player := initializePlayer(world)

	go render(world, &player)

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyCtrlQ {
				os.Exit(3)
			}

			moveVector := Vector{x: 0, y: 0}
			if ev.Ch == 119 { // w
				moveVector.y = -1
				movePlayer(world, &player, moveVector)
			} else if ev.Ch == 97 { // a
				moveVector.x = -1
				movePlayer(world, &player, moveVector)
			} else if ev.Ch == 115 { // s
				moveVector.y = 1
				movePlayer(world, &player, moveVector)
			} else if ev.Ch == 100 { // d
				moveVector.x = 1
				movePlayer(world, &player, moveVector)
			} else if ev.Ch == 0 { // space
				player.BuildWall(world)
			}

			checkAlive(&player)

			termbox.Flush()
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}

func initializeWorld() World {
	subWorlds := [WORLD_SIZE][WORLD_SIZE]SubWorld{}

	for i := 0; i < WORLD_SIZE; i++ {
		for j := 0; j < WORLD_SIZE; j++ {
			subWorlds[i][j] = initializeSubWorld()
		}
	}

	return World{subWorlds: subWorlds}
}

func initializeSubWorld() SubWorld {
	subWorld := SubWorld{}

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

func spawnSnakes(world *World) {
	for {
		go spawnSnake(world)
		time.Sleep(3000 * time.Millisecond)
	}
}

func runWorldElements(world *World) {
	//go spawnSnakes(world)
	go spawnCoinsInWorld(world)
}

func main() {
	rand.Seed(12345)

	firstBoot = true

	initializeRedisClient()

	world := initializeWorld()

	runWorldElements(&world)

	//startTerminalClient(&world)

	//val, err := redisClient.Get("key").Result()
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("key", val)

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
