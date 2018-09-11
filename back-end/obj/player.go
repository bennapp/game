package obj

import (
	"../typ"
)

const PLAYER = "player"

type Player struct {
	Object
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

func NewPlayer() *Player {
	return &Player{Object: newObject(PLAYER)}
}

func LoadPlayer() typ.Typical {
	return &Player{Object: loadObject(PLAYER)}
}
