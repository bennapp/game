package el

import (
	"../rc"
	"github.com/google/uuid"
	"math/rand"
)

const COIN = "coin"
const MAX_COIN_AMOUNT = 10

type Coin struct {
	Amount int
	Type   string
	Object
}

func (coin *Coin) String() string {
	return "C"
}

func (coin *Coin) Id() string {
	return coin.Object.Id.String()
}

func newCoin() *Coin {
	uuid, _ := uuid.NewUUID()
	amount := rand.Intn(MAX_COIN_AMOUNT)
	return &Coin{Type: COIN, Object: Object{Id: uuid}, Amount: amount}
}

func loadCoin() rc.Dbo {
	return &Coin{Type: COIN}
}
