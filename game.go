package main

import "github.com/nsf/termbox-go"
import (
	"fmt"
	"math/rand"
	"sync"
	"time"
	"os"
	"os/exec"
)

type Board [BOARD_SIZE][BOARD_SIZE]int
type GameBoard struct {
	board Board
	coinCount int
	mux sync.Mutex
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
const BOARD_SIZE = 36

func placePlayer(gameBoard *GameBoard) {
	placeElementRandomLocationWithLock(gameBoard, PLAYER)
}

func placeSnake(gameBoard *GameBoard) {
	placeElementRandomLocationWithLock(gameBoard, SNAKE)
}

func placeCoin(gameBoard *GameBoard) {
	placeElementRandomLocationWithLock(gameBoard, COIN)
}

func placeRock(gameBoard *GameBoard) {
	placeElementRandomLocationWithLock(gameBoard, ROCK)
}

func placeElementRandomLocationWithLock(gameBoard *GameBoard, element int) {
	gameBoard.mux.Lock()

	placeElementRandomLocation(gameBoard, element)

	gameBoard.mux.Unlock()
}

func placeElementRandomLocation(gameBoard *GameBoard, element int) {
	x := rand.Intn(BOARD_SIZE)
	y := rand.Intn(BOARD_SIZE)

	if gameBoard.board[x][y] == EMPTY {
		gameBoard.board[x][y] = element
	} else {
		placeElementRandomLocation(gameBoard, element)
	}
}

func findPlayer(gameBoard *GameBoard) Coord {
	return findElement(gameBoard, PLAYER)
}

func findSnake(gameBoard *GameBoard) Coord {
	return findElement(gameBoard, SNAKE)
}

func findElement(gameBoard *GameBoard, element int) Coord {
	x := -1
	y := -1

	for i := 0; i < BOARD_SIZE; i++ {
		for j := 0; j < BOARD_SIZE; j++ {
			if gameBoard.board[i][j] == element {
				x = i
				y = j
				break
			}
		}
	}

	return Coord{x: x, y: y}
}



func moveCharacter(gameBoard *GameBoard, coord Coord, vector Vector, element int) {
	gameBoard.mux.Lock()

	nextX := wrapAroundBoard(coord.x + vector.x)
	nextY := wrapAroundBoard(coord.y + vector.y)

	nextCoord := Coord{x: nextX, y: nextY }

	checkKillSnake(gameBoard, nextCoord)
	checkPickUpCoin(gameBoard, nextCoord, element)

	if !checkRock(gameBoard, nextCoord) {
		gameBoard.board[coord.x][coord.y] = 0
		gameBoard.board[nextX][nextY] = element
	}

	gameBoard.mux.Unlock()
}

//No lock
func checkKillSnake(gameBoard *GameBoard, coord Coord) {
	if gameBoard.board[coord.x][coord.y] == SNAKE {
		fmt.Println("YUM!!!")
	}
}

//No lock
func checkPickUpCoin(gameBoard *GameBoard, coord Coord, element int) {
	if element != PLAYER {
		return
	}

	if gameBoard.board[coord.x][coord.y] == COIN {
		gameBoard.coinCount++
	}
}

//No lock
func checkRock(gameBoard *GameBoard, coord Coord) bool {
	return gameBoard.board[coord.x][coord.y] == ROCK
}

func wrapAroundBoard(n int) int {
	if n == -1 {
		return BOARD_SIZE - 1
	}
	if n == BOARD_SIZE {
		return 0
	}

	return n
}

func printBoard(x Board) {
	for _, i := range x {
		for _, j := range i {
			fmt.Printf("%d ", j)
		}
		fmt.Println()
	}
	fmt.Println()
}

func printStat(gameBoard *GameBoard) {
	fmt.Printf("Coin: %d", gameBoard.coinCount)
	fmt.Println()
}

func showGame(gameBoard *GameBoard) {
	for {
		printBoard(gameBoard.board)
		printStat(gameBoard)
		time.Sleep(250 * time.Millisecond)
		clearScreen()
	}
}

func clearScreen() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
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

func snakeWalk(gameBoard *GameBoard) {
	for {
		snakeLocation := findSnake(gameBoard)
		if snakeLocation.x == -1 && snakeLocation.y == -1 {
			return
		}

		playerLocation := findPlayer(gameBoard)

		diffX := playerLocation.x - snakeLocation.x
		diffY := playerLocation.y - snakeLocation.y

		moveVector := Vector{x: 0, y: 0}
		if abs(diffX) > abs(diffY) {
			moveVector.x = convertToOneMove(diffX)
		} else {
			moveVector.y = convertToOneMove(diffY)
		}

		moveCharacter(gameBoard, snakeLocation, moveVector, SNAKE)
		time.Sleep(1000 * time.Millisecond)
	}
}

func dropGold(gameBoard *GameBoard) {
	for {
		placeCoin(gameBoard)
		time.Sleep(4000 * time.Millisecond)
	}
}

func startTerminalClient(gameBoard *GameBoard) {
	go showGame(gameBoard)

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyCtrlQ {
				os.Exit(3)
			}

			moveVector := Vector{x: 0, y: 0}
			if ev.Ch == 119 { // w
				moveVector = Vector{x: -1, y: 0}
			} else if ev.Ch == 97 { // a
				moveVector = Vector{x: 0, y: -1}
			} else if ev.Ch == 115 { // s
				moveVector = Vector{x: 1, y: 0}
			} else if ev.Ch == 100 { // d
				moveVector = Vector{x: 0, y: 1}
			}

			playerLocation := findPlayer(gameBoard)
			moveCharacter(gameBoard, playerLocation, moveVector, PLAYER)

			termbox.Flush()
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	board := Board{}
	gameBoard := GameBoard{board: board}

	placePlayer(&gameBoard)
	placeSnake(&gameBoard)

	for i:=0; i < 10; i++ {
		placeRock(&gameBoard)
	}

	go snakeWalk(&gameBoard)
	go dropGold(&gameBoard)

	startTerminalClient(&gameBoard)
}
