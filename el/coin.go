package el

import (
	"../rc"
	"fmt"
	"strconv"
	"strings"
)

const COIN = "coin"

type Coin struct {
	Location
	Amount int
	id     int
}

func (coin *Coin) String() string {
	return "C"
}

func (coin *Coin) Type() string {
	return COIN
}

func (coin *Coin) Id() int {
	return coin.id
}

func (coin *Coin) Key() string {
	return rc.GenerateKey(COIN, coin.id)
}

func (coin *Coin) Serialize() string {
	return fmt.Sprintf("amount:%v", coin.Amount)
}

func (coin *Coin) Deserialize(key string, values string) {
	amountString := strings.Split(values, "amount:")[1]
	amount, _ := strconv.Atoi(amountString)

	id, _ := rc.SplitKey(key)

	coin.id, _ = strconv.Atoi(id)
	coin.Amount = amount
}

func newCoinDbo(id int) rc.Dbo {
	return &Coin{id: id}
}
