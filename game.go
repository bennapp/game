package main

import "github.com/nsf/termbox-go"
import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

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
	element Element
}
type Element struct {
	subWorldCoord Coord
	gridCoord     Coord
}
type Player struct {
	Element
	mux       sync.Mutex
	coinCount int
	alive     bool
}
type Snake struct {
	Element
}
type Coin struct {
	Element
}
type Rock struct {
	Element
}

const GRID_SIZE = 8
const WORLD_SIZE = 4

// creates a player
// returns a pair of Coord of World, SubWorld
func initializePlayer(world *World) Player {
	x, y := randomPair(WORLD_SIZE)
	subWorld := world.subWorlds[x][y]

	player := Player{alive: true}
	subWorldCoord := Coord{x: x, y: y}
	gridCoord := placeElementRandomLocationWithLock(&subWorld, player.cell)

	player.Character.cell = playerCell
	player.subWorldCoord = subWorldCoord
	player.gridCoord = gridCoord

	return player
}

// creates a snakeCell
// returns a pair of Coord of World, SubWorld
func placeSnake(world *World) (Coord, Coord) {
	x, y := randomPair(WORLD_SIZE)
	subWorld := world.subWorlds[x][y]

	return Coord{x: x, y: y}, placeElementRandomLocationWithLock(&subWorld, snakeCell)
}

// creates a Coin
// returns a pair of Coord of SubWorld
func placeCoin(subWorld *SubWorld) Coord {
	return placeElementRandomLocationWithLock(subWorld, coinCell)
}

// creates a Rock
// returns a pair of Coord of SubWorld
func placeRock(subWorld *SubWorld) Coord {
	return placeElementRandomLocationWithLock(subWorld, rockCell)
}

func placeElementRandomLocationWithLock(subWorld *SubWorld, element Cell) Coord {
	x, y := randomPair(GRID_SIZE)
	coord := Coord{x: x, y: y}

	if subWorld.grid[x][y].code == emptyCell.code {
		cell := &subWorld.grid[x][y]

		cell.mux.Lock()
		cell.code = element.code
		cell.mux.Unlock()
	} else {
		coord = placeElementRandomLocationWithLock(subWorld, element)
	}

	return coord
}

func movePlayer(world *World, player *Player, vector Vector) {
	player.subWorldCoord, player.gridCoord = moveCharacter(world, player.subWorldCoord, player.gridCoord, vector, &player.Character)
}

func moveCharacter(world *World, subWorldCoord Coord, coord Coord, vector Vector, character *Character) (Coord, Coord) {
	subWorld := &world.subWorlds[subWorldCoord.x][subWorldCoord.y]

	nextSubWorldCoord, nextCoord := subWorldMove(subWorldCoord, coord, vector)

	nextSubWorld := &world.subWorlds[nextSubWorldCoord.x][nextSubWorldCoord.y]

	nextElement := getElement(nextSubWorld, nextCoord)
	interactFunc := elementInteractFuncMap[character.cell.code][nextElement.code]
	character.Interact(nextElement)

	override := false
	if interactFunc != nil {
		override = interactFunc(world)
	} else {
		override = true
	}

	if override {
		prevCell := &subWorld.grid[coord.x][coord.y]

		prevCell.mux.Lock()
		prevCell.code = emptyCell.code
		prevCell.mux.Unlock()

		nextCell := &nextSubWorld.grid[nextCoord.x][nextCoord.y]
		nextCell.mux.Lock()
		nextCell.code = character.cell.code
		nextCell.mux.Unlock()
	} else {
		nextSubWorldCoord = subWorldCoord
		nextCoord = coord
	}
	return nextSubWorldCoord, nextCoord
}

func getElement(subWorld *SubWorld, coord Coord) Cell {
	return subWorld.grid[coord.x][coord.y]
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

// return true if element can replace snakeCell.code
func interactWithSnake(world *World) bool {
	fmt.Println("Argh run!")
	return true
}

// return true if element can replace rock.code
func interactWithCoin(world *World) bool {
	world.mux.Lock()
	world.coinCount++
	world.mux.Unlock()

	return true
}

// return true if element can replace rock.code
func interactWithRock(world *World) bool {
	return false
}

func killPlayer(world *World) bool {
	world.mux.Lock()
	world.alive = false
	world.mux.Unlock()

	return true
}

func printWorld(world *World) {
	for wy := 0; wy < WORLD_SIZE; wy++ {
		for gy := 0; gy < GRID_SIZE; gy++ {
			for wx := 0; wx < WORLD_SIZE; wx++ {
				for gx := 0; gx < GRID_SIZE; gx++ {
					fmt.Printf("%v ", cellDisplayLookup[world.subWorlds[wx][wy].grid[gx][gy].code])
				}
				fmt.Printf("|")
			}
			fmt.Println()
		}
		fmt.Println(strings.Repeat(" -", (WORLD_SIZE*GRID_SIZE)+2))
	}
}

func printStat(world *World) {
	fmt.Printf("Coin: %d", world.coinCount)
	fmt.Println()

	if !world.alive {
		fmt.Println("You Died")
		os.Exit(0)
	}
}

func render(world *World) {
	for {
		clearScreen()
		printWorld(world)
		printStat(world)
		time.Sleep(250 * time.Millisecond)
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
	return findElement(subWorld, playerCell)
}

func findElement(subWorld *SubWorld, element Cell) Coord {
	x := -1
	y := -1

	for i := 0; i < GRID_SIZE; i++ {
		for j := 0; j < GRID_SIZE; j++ {
			if subWorld.grid[i][j].code == element.code {
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

		subWorldCoord, snakeCoord = moveCharacter(world, subWorldCoord, snakeCoord, moveVector, snakeCell)
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

	player := initializePlayer(world)

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

			movePlayer(world, &player, moveVector)

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

	for i := 0; i < GRID_SIZE; i++ {
		for j := 0; j < GRID_SIZE; j++ {
			subWorld.grid[i][j].code = emptyCell.code
		}
	}

	for i := 0; i < 10; i++ {
		placeRock(&subWorld)
	}

	return subWorld
}

func spawnSnakes(world *World) {
	for {
		go snakeWalk(world)
		time.Sleep(2000 * time.Millisecond)
	}
}

func initializeWorldElements(world *World) {
	go spawnSnakes(world)
	go spawnGoldInWorld(world)
}

func main() {
	rand.Seed(12345)

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
