package obj

import (
	"../gs"
	"../typ"
)

const PLAYER = "player"

type Player struct {
	Object
	Location  gs.Coord
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
	return &Player{Object: newObject(PLAYER), Location: location}
}

func NewPlayer() *Player {
	coord := gs.NewRandomCoord()
	return &Player{Object: newObject(PLAYER), Location: coord}
}

func LoadPlayer() typ.Typical {
	return &Player{Object: loadObject(PLAYER)}
}
