package el

import (
	"../rc"
	"github.com/google/uuid"
	"sync"
)

const PLAYER = "player"

type Player struct {
	Location
	mux       sync.Mutex
	CoinCount int
	Alive     bool
	Hp        int
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

func newPlayerDbo(isNew bool) rc.Dbo {
	player := &Player{Type: PLAYER}

	if isNew {
		player.Id = uuid.New()
	}

	return player
}

func (player *Player) Mux() *sync.Mutex {
	return &player.mux
}
