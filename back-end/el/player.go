package el

import (
	"../rc"
	"sync"
	"github.com/google/uuid"
)

const PLAYER = "player"

type Player struct {
	Location
	mux       sync.Mutex
	CoinCount int
	Alive     bool
	Hp        int
	Avatar    string
	Id        uuid.UUID
	Type      string
}

func (player Player) String() string {
	return "P"
}

func (player *Player) Key() string {
	return rc.GenerateKey(player.Id)
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

func (player *Player) Load() rc.Dbo {
	return &Player{Type: PLAYER}
}

func newPlayerDbo(id uuid.UUID) rc.Dbo {
	return &Player{Id: id, Type: PLAYER}
}

func (player *Player) Mux() *sync.Mutex {
	return &player.mux
}
