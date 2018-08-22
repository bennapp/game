package el

import (
	"fmt"
	"sync"
)

type Player struct {
	Element
	mux       sync.Mutex
	coinCount int
	alive     bool
	hp        int
	id        int
}

func (p Player) String() string {
	return fmt.Sprintf("%v", p.id)
}
func (player *Player) Interact(element interface{}) bool {
	switch v := element.(type) {
	case Coin:
		player.IncCoinCount(v.amount)
		v.Destroy()
		return true
	case Empty:
		return true
	default:
		return false
	}

	return false
}
func (player *Player) Kill() {
	player.mux.Lock()
	player.alive = false
	player.mux.Unlock()
}
func (player *Player) IncCoinCount(amount int) {
	player.mux.Lock()
	player.coinCount += amount
	storePlayer(player)
	player.mux.Unlock()
}
func (player *Player) decreaseHp(damage int) {
	player.mux.Lock()
	player.hp -= damage
	player.mux.Unlock()

	if player.hp < 0 {
		player.Kill()
	}
}
func (player *Player) Id() string {
	return fmt.Sprintf("player:%v", player.id)
}
func (player *Player) Val() string {
	// Bug fix, use dashes because cords use commas. FIXME: use commas for all attr delimiters
	return fmt.Sprintf("coinCount:%v-alive:%v-hp:%v-subWorldCoord:%v-gridCoord:%v",
		player.coinCount,
		player.alive,
		player.hp,
		player.SubWorldCoord.Key(),
		player.GridCoord.Key(),
	)
}

func NewPlayer(id int, coinCount int, alive bool, hp int) Player {
	return Player{id: id, coinCount: coinCount, alive: alive, hp: hp}
}

func (player *Player) Mux() *sync.Mutex {
	return &player.mux
}

func (player *Player) CoinCount() int {
	return player.coinCount
}

func (player *Player) Alive() bool {
	return player.alive
}

func (player *Player) Hp() int {
	return player.hp
}
