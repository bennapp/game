package obj

import (
	"../gs"
	"../mov"
	"../typ"
	"../velocity"
)

const PLAYER = "player"

type Player struct {
	Object
	mov.Mover
	CoinCount int
	Alive     bool
	Hp        int
	Avatar    string
}

func (player *Player) Kill() {
	player.Alive = false
}

func (player *Player) IncCoinCount(amount int) {
	player.CoinCount += amount
}

func (player *Player) DecreaseHp(damage int) {
	player.Hp -= damage

	if player.Hp < 0 {
		player.Kill()
	}
}

func NewPlayerAt(location gs.Coord) *Player {
	player := NewPlayer()
	player.SetLocation(location)
	return player
}

func NewPlayer() *Player {
	return &Player{Object: newObject(PLAYER), Hp: 10, Alive: true}
}

func LoadPlayer() typ.Typical {
	return &Player{}
}

func (player *Player) GetVelocity() float64 {
	return velocity.Constants(PLAYER)
}

func IsPlayer(entity interface{}) bool {
	entityIsPlayer := false

	switch entity.(type) {
	case *Player:
		entityIsPlayer = true
	}

	return entityIsPlayer
}
