package el

import (
	"../rc"
)

const ROCK = "rock"

type Rock struct {
	Type string
}

func (rock *Rock) String() string {
	return "R"
}

func (rock *Rock) Id() int {
	return -1
}

func (rock *Rock) Key() string {
	return rc.GenerateKey(ROCK, -1)
}

func NewRock() *Rock {
	return &Rock{}
}

func newRockDbo(id int) rc.Dbo {
	return &Rock{Type: ROCK}
}
