package main

import "github.com/nsf/termbox-go"
import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sync"
	"time"
)

type World struct {
	subWorlds [WORLD_SIZE][WORLD_SIZE]SubWorld
	coinCount int
	mux       sync.Mutex // remove after coin count refactor
}
type Grid [GRID_SIZE][GRID_SIZE]int
type SubWorld struct {
	grid Grid
	mux  sync.Mutex
}
type Coord struct {
	x int
	y int
}
type Vector struct {
	x int
	y int
}

const EMPTY = 0
const PLAYER = 1
const SNAKE = 2
const COIN = 3
const ROCK = 4
const GRID_SIZE = 8
const WORLD_SIZE = 2

// creates a player
// returns a pair of Coord of World, SubWorld
func initializePlayer(world *World) (Coord, Coord) {
	x, y := randomPair(WORLD_SIZE)
	subWorld := world.subWorlds[x][y]

	return Coord{x: x, y: y}, placeElementRandomLocationWithLock(&subWorld, PLAYER)
}

// creates a snake
// returns a pair of Coord of World, SubWorld
func placeSnake(world *World) (Coord, Coord) {
	x, y := randomPair(WORLD_SIZE)
	subWorld := world.subWorlds[x][y]

	return Coord{x: x, y: y}, placeElementRandomLocationWithLock(&subWorld, SNAKE)
}

// creates a Coin
// returns a pair of Coord of SubWorld
func placeCoin(subWorld *SubWorld) Coord {
	return placeElementRandomLocationWithLock(subWorld, COIN)
}

// creates a Rock
// returns a pair of Coord of SubWorld
func placeRock(subWorld *SubWorld) Coord {
	return placeElementRandomLocationWithLock(subWorld, ROCK)
}

// creates a Rock
// returns a pair of Coord of SubWorld
func placeElementRandomLocationWithLock(subWorld *SubWorld, element int) Coord {
	subWorld.mux.Lock()

	coord := placeElementRandomLocation(subWorld, element)

	subWorld.mux.Unlock()

	return coord
}

func placeElementRandomLocation(subWorld *SubWorld, element int) Coord {
	x, y := randomPair(GRID_SIZE)
	coord := Coord{x: x, y: y}

	if subWorld.grid[x][y] == EMPTY {
		subWorld.grid[x][y] = element
	} else {
		coord = placeElementRandomLocation(subWorld, element)
	}

	return coord
}

func moveCharacter(world *World, subWorldCoord Coord, coord Coord, vector Vector, element int) (Coord, Coord) {
	prevSubWorld := &world.subWorlds[subWorldCoord.x][subWorldCoord.y]

	subWorldCoord, nextCoord := subWorldMove(subWorldCoord, coord, vector)

	nextSubWorld := &world.subWorlds[subWorldCoord.x][subWorldCoord.y]

	checkKillSnake(prevSubWorld, nextCoord)
	checkPickUpCoin(world, prevSubWorld, nextCoord, element)

	if !checkRock(prevSubWorld, nextCoord) {
		prevSubWorld.mux.Lock()
		prevSubWorld.grid[coord.x][coord.y] = EMPTY
		prevSubWorld.mux.Unlock()

		nextSubWorld.mux.Lock()
		nextSubWorld.grid[nextCoord.x][nextCoord.y] = element
		nextSubWorld.mux.Unlock()
	}

	return subWorldCoord, nextCoord
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

//No lock
func checkKillSnake(subWorld *SubWorld, coord Coord) {
	if subWorld.grid[coord.x][coord.y] == SNAKE {
		fmt.Println("YUM!!!")
	}
}

//No lock
func checkPickUpCoin(world *World, subWorld *SubWorld, coord Coord, element int) {
	if element != PLAYER {
		return
	}

	if subWorld.grid[coord.x][coord.y] == COIN {
		world.mux.Lock()
		world.coinCount++
		world.mux.Unlock()
	}
}

//No lock
func checkRock(subWorld *SubWorld, coord Coord) bool {
	return subWorld.grid[coord.x][coord.y] == ROCK
}

func printWorld(world *World) {
	for wy := 0; wy < WORLD_SIZE; wy++ {
		for gy := 0; gy < GRID_SIZE; gy++ {
			for wx := 0; wx < WORLD_SIZE; wx++ {
				for gx := 0; gx < GRID_SIZE; gx++ {
					fmt.Printf("%d ", world.subWorlds[wx][wy].grid[gx][gy])
				}
				fmt.Printf(" ")
			}
			fmt.Println()
		}
		fmt.Println()
	}
}

func printStat(world *World) {
	fmt.Printf("Coin: %d", world.coinCount)
	fmt.Println()
}

func render(world *World) {
	for {
		printWorld(world)
		printStat(world)
		time.Sleep(1000 * time.Millisecond)
		clearScreen()
	}
}

func clearScreen() {
	cmd := exec.Command("cmd", "/c", "cls || clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

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

func findPlayer(subWorld *SubWorld) Coord {
	return findElement(subWorld, PLAYER)
}

func findElement(subWorld *SubWorld, element int) Coord {
	x := -1
	y := -1

	for i := 0; i < GRID_SIZE; i++ {
		for j := 0; j < GRID_SIZE; j++ {
			if subWorld.grid[i][j] == element {
				x = i
				y = j
				break
			}
		}
	}

	return Coord{x: x, y: y}
}

func snakeWalk(world *World) {
	subWorldCoord, snakeCoord := placeSnake(world)

	for {
		if snakeCoord.x == -1 && snakeCoord.y == -1 {
			return
		}

		playerLocation := findPlayer(&world.subWorlds[subWorldCoord.x][subWorldCoord.y])

		moveVector := Vector{x: 0, y: 0}

		if isFound(playerLocation) {
			diffX := playerLocation.x - snakeCoord.x
			diffY := playerLocation.y - snakeCoord.y

			if abs(diffX) > abs(diffY) {
				moveVector.x = convertToOneMove(diffX)
			} else {
				moveVector.y = convertToOneMove(diffY)
			}
		} else {
			moveVector = randomVector()
		}

		subWorldCoord, snakeCoord = moveCharacter(world, subWorldCoord, snakeCoord, moveVector, SNAKE)
		time.Sleep(1000 * time.Millisecond)
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

func spawnGoldInWorld(world *World) {
	for {
		randomSubWorld := randomSubWorld(world)
		spawnGoldInSubWorld(randomSubWorld)
		time.Sleep(4000 * time.Millisecond)
	}
}

func spawnGoldInSubWorld(subWorld *SubWorld) {
	placeCoin(subWorld)
}

func startTerminalClient(world *World) {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	subWorldCoord, playerCoord := initializePlayer(world)

	go render(world)

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyCtrlQ {
				os.Exit(3)
			}

			moveVector := Vector{x: 0, y: 0}
			if ev.Ch == 119 { // w
				moveVector.y = -1
			} else if ev.Ch == 97 { // a
				moveVector.x = -1
			} else if ev.Ch == 115 { // s
				moveVector.y = 1
			} else if ev.Ch == 100 { // d
				moveVector.x = 1
			}

			subWorldCoord, playerCoord = moveCharacter(world, subWorldCoord, playerCoord, moveVector, PLAYER)
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
			subWorlds[i][j] = initializeSubworld()
		}
	}

	return World{subWorlds: subWorlds}
}

func initializeSubworld() SubWorld {
	subWorld := SubWorld{}

	for i := 0; i < 10; i++ {
		placeRock(&subWorld)
	}

	return subWorld
}

func initializeWorldElements(world *World) {
	go snakeWalk(world)
	go spawnGoldInWorld(world)
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	world := initializeWorld()

	initializeWorldElements(&world)

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
