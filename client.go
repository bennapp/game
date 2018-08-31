package main

import (
	"./el"
	"./gs"
	"./wo"
	"fmt"
	"github.com/nsf/termbox-go"
	"os"
	"time"
)

func clearScreen() {
	//cmd := exec.Command("cmd", "/c", "cls || clear")
	//cmd.Stdout = os.Stdout
	//cmd.Run()

	print("\033[H\033[2J")
}

func printWorld(player *el.Player) {
	v := gs.NewVector(-5, -5)
	visionDistance := 11

	for i := 0; i < visionDistance; i++ {
		for j := 0; j < visionDistance; j++ {
			element, valid := wo.NextElement(player.GridCoord, v)
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

func render(player *el.Player) {
	for {
		clearScreen()
		printWorld(player)
		fmt.Println()
		printStat(player)

		time.Sleep(100 * time.Millisecond)
	}
}

func startTerminalClient(id int, char string) {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	player := wo.LoadPlayer(id)
	player.Avatar = char

	go render(player)

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyCtrlQ {
				wo.PlayerLogout(player)
				wo.Close()
				termbox.Close()
				clearScreen()
				os.Exit(3)
			}

			moveVector := gs.NewVector(0, 0)
			if ev.Ch == 119 { // w
				moveVector.Y = -1
				wo.MovePlayer(player, moveVector)
			} else if ev.Ch == 97 { // a
				moveVector.X = -1
				wo.MovePlayer(player, moveVector)
			} else if ev.Ch == 115 { // s
				moveVector.Y = 1
				wo.MovePlayer(player, moveVector)
			} else if ev.Ch == 100 { // d
				moveVector.X = 1
				wo.MovePlayer(player, moveVector)
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

func main() {
	wo.Init()

	id := 386
	char := "M"

	startTerminalClient(id, char)
}
