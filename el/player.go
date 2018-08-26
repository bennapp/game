package el

import (
	"../gs"
	"../rc"
	"fmt"
	"strconv"
	"strings"
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
	id        int
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
	return rc.GenerateKey(PLAYER, player.id)
}

func (player *Player) Serialize() string {
	// Bug fix, use dashes because cords use commas. FIXME: use commas for all attr delimiters
	// TODO - consider using json
	return fmt.Sprintf("CoinCount:%v-Alive:%v-Hp:%v-subWorldCoord:%v-gridCoord:%v",
		player.CoinCount,
		player.Alive,
		player.Hp,
		player.SubWorldCoord.Key(),
		player.GridCoord.Key(),
	)
}

func (player *Player) Deserialize(key string, values string) {
	keyValues := strings.Split(values, "-")

	coinCountString := strings.Split(keyValues[0], "CoinCount:")[1]
	aliveString := strings.Split(keyValues[1], "Alive:")[1]
	hpString := strings.Split(keyValues[2], "Hp:")[1]

	subWorldCoordString := strings.Split(keyValues[3], "subWorldCoord:")[1]
	subWorldCoordStringX := strings.Split(subWorldCoordString, ",")[0]
	subWorldCoordX, _ := strconv.Atoi(subWorldCoordStringX)
	subWorldCoordStringY := strings.Split(subWorldCoordString, ",")[1]
	subWorldCoordY, _ := strconv.Atoi(subWorldCoordStringY)

	gridCoordString := strings.Split(keyValues[4], "gridCoord:")[1]
	gridCoordStringX := strings.Split(gridCoordString, ",")[0]
	gridCoordX, _ := strconv.Atoi(gridCoordStringX)
	gridCoordStringY := strings.Split(gridCoordString, ",")[1]
	gridCoordY, _ := strconv.Atoi(gridCoordStringY)

	coinCount, _ := strconv.Atoi(coinCountString)
	hp, _ := strconv.Atoi(hpString)
	alive := aliveString == "true"

	_, id := rc.SplitKey(key)

	player.id, _ = strconv.Atoi(id)
	player.CoinCount = coinCount
	player.Alive = alive
	player.Hp = hp
	player.SubWorldCoord = gs.NewCoord(subWorldCoordX, subWorldCoordY)
	player.GridCoord = gs.NewCoord(gridCoordX, gridCoordY)
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

func NewPlayer(id int, coinCount int, alive bool, hp int) Player {
	return Player{id: id, CoinCount: coinCount, Alive: alive, Hp: hp}
}

func newPlayerDbo(id int) rc.Dbo {
	return &Player{id: id}
}

func (player *Player) Mux() *sync.Mutex {
	return &player.mux
}
