package main

import "github.com/nsf/termbox-go"
import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Board [BOARD_SIZE][BOARD_SIZE]int
type GameBoard struct {
	board Board
	mux   sync.Mutex
}
type Coord struct {
	x int
	y int
}
type Vector struct {
	x int
	y int
}

const PLAYER = 1
const SNAKE = 2
const BOARD_SIZE = 8

func placePlayer(gameBoard *GameBoard) {
	placeElementRandomLocation(gameBoard, PLAYER)
}

func placeSnake(gameBoard *GameBoard) {
	placeElementRandomLocation(gameBoard, SNAKE)
}

func placeElementRandomLocation(gameBoard *GameBoard, element int) {
	gameBoard.mux.Lock()

	x := rand.Intn(BOARD_SIZE)
	y := rand.Intn(BOARD_SIZE)

	gameBoard.board[x][y] = element

	gameBoard.mux.Unlock()
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

	gameBoard.board[coord.x][coord.y] = 0

	nextX := wrap(coord.x + vector.x)
	nextY := wrap(coord.y + vector.y)

	snakeLocation := findSnake(gameBoard)

	if snakeLocation.x == nextX && snakeLocation.y == nextY {
		fmt.Println("YUM!!!")
	}

	gameBoard.board[nextX][nextY] = element

	gameBoard.mux.Unlock()
}

func wrap(n int) int {
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

func showGame(gameBoard *GameBoard) {
	for {
		printBoard(gameBoard.board)
		time.Sleep(250 * time.Millisecond)
		clearScreen()
	}
}

func clearScreen() {
	print("\033[H\033[2J")
}

func abs(n int) int {
	if n < 0 {
		return -n
	} else {
		return n
	}
}

func normalize(n int) int {
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
			moveVector.x = normalize(diffX)
		} else {
			moveVector.y = normalize(diffY)
		}

		moveCharacter(gameBoard, snakeLocation, moveVector, SNAKE)
		time.Sleep(1000 * time.Millisecond)
	}
}

func startTerminalClient(gameBoard *GameBoard) {
	go showGame(gameBoard)

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyCtrlQ {
				break loop
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

			//termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			//draw_keyboard()
			//dispatch_press(&ev)
			//pretty_print_press(&ev)
			termbox.Flush()
		case termbox.EventResize:
			//termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			//draw_keyboard()
			//pretty_print_resize(&ev)
			//termbox.Flush()
		case termbox.EventMouse:
			//termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			//draw_keyboard()
			//pretty_print_mouse(&ev)
			//termbox.Flush()
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

	go snakeWalk(&gameBoard)

	startTerminalClient(&gameBoard)
}
