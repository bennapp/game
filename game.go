package main

import "github.com/nsf/termbox-go"
import (
	"fmt"
	"time"
	"math/rand"
	"sync"
	)

type Board [8][8]int
type GameBoard struct {
	board Board
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

const PLAYER = 1
const SNAKE = 2

func placePlayer(gameBoard *GameBoard) {
	placeElementRandomLocation(gameBoard, PLAYER)
}

func placeSnake(gameBoard *GameBoard) {
	placeElementRandomLocation(gameBoard, SNAKE)
}

func placeElementRandomLocation(gameBoard *GameBoard, element int) {
	gameBoard.mux.Lock()

	x := rand.Intn(8)
	y := rand.Intn(8)

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

	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
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

	if (snakeLocation.x == nextX && snakeLocation.y == nextY) {
		fmt.Println("YUM!!!")
	}

	gameBoard.board[nextX][nextY] = element

	gameBoard.mux.Unlock()
}

func wrap(n int) int {
	if n == -1 {
		return 8 - 1
	}
	if n == 8 {
		return 0
	}

	return n
}

func print(x Board) {
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
		print(gameBoard.board)
		time.Sleep(250 * time.Millisecond)
	}
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
		if (snakeLocation.x == -1 && snakeLocation.y == -1) {
			return
		}

		playerLocation := findPlayer(gameBoard)

		diffX := playerLocation.x - snakeLocation.x
		diffY := playerLocation.y - snakeLocation.y

		moveVector := Vector{x: 0 , y: 0}
		if abs(diffX) > abs(diffY) {
			moveVector.x = normalize(diffX)
		} else {
			moveVector.y = normalize(diffY)
		}

		moveCharacter(gameBoard, snakeLocation, moveVector, SNAKE)
		time.Sleep(1000 * time.Millisecond)
	}
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	board := Board{}
	gameBoard := GameBoard{board: board}

	placePlayer(&gameBoard)
	placeSnake(&gameBoard)

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	go showGame(&gameBoard)
	go snakeWalk(&gameBoard)
loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyCtrlQ {
				break loop
			}

			moveVector := Vector{x:0 , y:0}
			if ev.Ch == 119 { // w
				moveVector = Vector{x: -1, y: 0}
			} else if ev.Ch == 97 { // a
				moveVector = Vector{x: 0, y: -1}
			} else if ev.Ch == 115 { // s
				moveVector = Vector{x: 1, y: 0}
			} else if ev.Ch == 100 { // d
				moveVector = Vector{x: 0, y: 1}
			}

			playerLocation := findPlayer(&gameBoard)
			moveCharacter(&gameBoard, playerLocation, moveVector, PLAYER)

			//print(board)

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
