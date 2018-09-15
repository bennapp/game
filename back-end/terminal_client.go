package main

import (
	"./dbs"
	"./gs"
	"./obj"
	"./wo"
	"fmt"
	"github.com/nsf/termbox-go"
	"os"
	"os/exec"
	"time"
)

func clearScreen() {
	cmd := exec.Command("cmd", "/c", "cls || clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	//print("\033[H\033[2J")
}

func printWorld(player *obj.Player) {
	visionDistance := gs.LOADED_VISION_DISTANCE
	halfVisionDistance := visionDistance / 2
	v := gs.NewVector(-halfVisionDistance, -halfVisionDistance)

	for i := 0; i < visionDistance; i++ {
		for j := 0; j < visionDistance; j++ {
			coord := player.GetLocation().AddVector(v)
			cell := dbs.LoadCell(coord)
			fmt.Print(cell.DebugString())
			v.X += 1
		}
		v.X = -halfVisionDistance
		fmt.Println()
		v.Y += 1
	}
}

func printStat(player *obj.Player) {
	fmt.Printf("Coin: %d", player.CoinCount)
	fmt.Println()
	fmt.Printf("HP: %d", player.Hp)
	fmt.Println()
}

func checkAlive(player *obj.Player) {
	if !player.Alive {
		fmt.Println("You Died")
		os.Exit(0)
	}
}

func render(player *obj.Player) {
	for {
		clearScreen()
		printWorld(player)
		fmt.Println()
		printStat(player)

		time.Sleep(100 * time.Millisecond)
	}
}

func startTerminalClient() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	player := wo.CreatePlayer()
	go render(player)

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyCtrlQ {
				wo.DeletePlayer(player)
				termbox.Close()
				clearScreen()
				os.Exit(3)
			}

			if ev.Ch == 0 { // space
				//	player.BuildWall(gs)
			} else {
				moveVector := gs.NewVector(0, 0)
				if ev.Ch == 119 { // w
					moveVector.Y = -1
					wo.MoveObject(player, moveVector)
				} else if ev.Ch == 97 { // a
					moveVector.X = -1
					wo.MoveObject(player, moveVector)
				} else if ev.Ch == 115 { // s
					moveVector.Y = 1
					wo.MoveObject(player, moveVector)
				} else if ev.Ch == 100 { // d
					moveVector.X = 1
					wo.MoveObject(player, moveVector)
				}
			}

			checkAlive(player)

			termbox.Flush()
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}

func main() {
	startTerminalClient()
}
