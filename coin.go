package main

import (
	"./wo"
	"fmt"
	"github.com/nsf/termbox-go"
	"os"
)

func main() {
	fmt.Printf("coin.go: running World Elements\n")

	wo.Init()
	go wo.SpawnCoinsInWorld()

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyCtrlQ {
				wo.Close()
				termbox.Close()
				os.Exit(3)
			}
			termbox.Flush()
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}
