package main

import "github.com/nsf/termbox-go"
import (
	"fmt"
	"time"
	"math/rand"
)

type Board [8][8]int
type Coord struct {
	x int
	y int
}
type Vector struct {
	x int
	y int
}

func placeCharacter(board Board) Board {
	board = placeElementRandomLocation(board, 1)

	return board
}

func placeSnake(board Board) Board {
	board = placeElementRandomLocation(board, 2)

	return board
}

func placeElementRandomLocation(board Board, element int) Board {
	x := rand.Intn(8)
	y := rand.Intn(8)

	board[x][y] = element

	return board
}

func findCharacter(board Board) Coord {
	x := 0
	y := 0

	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if board[i][j] == 1 {
				x = i
				y = j
				break
			}
		}
	}

	return Coord{x: x, y: y}
}

func moveCharacter(board Board, coord Coord, vector Vector) Board {
	board[coord.x][coord.y] = 0

	nextX := wrap(coord.x + vector.x)
	nextY := wrap(coord.y + vector.y)

	fmt.Println(nextX)

	board[nextX][nextY] = 1

	return board
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

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	board := Board{}

	board = placeCharacter(board)
	board = placeSnake(board)

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

			playerLocation := findCharacter(board)
			board = moveCharacter(board, playerLocation, moveVector)
			print(board)

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