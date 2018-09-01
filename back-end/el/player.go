package el

import (
	"../rc"
	"sync"
)

const PLAYER = "player"

type Player struct {
	Location
	mux       sync.Mutex
	CoinCount int
	Alive     bool
	Hp        int
	Avatar    string //TODO - change this to limit to 1 character; also this is not saved
	Id        int
}

func (player Player) String() string {
	if player.Avatar == "" {
		player.Avatar = "P"
	}

	return player.Avatar
}

func (player Player) Type() string {
	return PLAYER
}

func (player *Player) Key() string {
	return rc.GenerateKey(PLAYER, player.Id)
}

func (player *Player) Interact(element rc.Dbo) bool {
	switch v := element.(type) {
	case *Coin:
		player.IncCoinCount(v.Amount)
		//v.Destroy()
		return true
	case *Empty:
		return true
	default:
		return false
	}

	return false
}

func (player *Player) Kill() {
	player.mux.Lock()
	player.Alive = false
	player.mux.Unlock()
}

func (player *Player) IncCoinCount(amount int) {
	player.mux.Lock()
	player.CoinCount += amount
	player.mux.Unlock()
}

func (player *Player) DecreaseHp(damage int) {
	player.mux.Lock()
	player.Hp -= damage
	player.mux.Unlock()

	if player.Hp < 0 {
		player.Kill()
	}
}

func newPlayerDbo(id int) rc.Dbo {
	return &Player{Id: id}
}

func (player *Player) Mux() *sync.Mutex {
	return &player.mux
}
