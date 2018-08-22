package el

import "fmt"

// COIN
type Coin struct {
	Element
	amount int
	id     int
}

func (c Coin) String() string {
	return "C"
}
func (coin *Coin) Id() string {
	return fmt.Sprintf("coin:%v", coin.id)
}
func (coin *Coin) Val() string {
	return fmt.Sprintf("amount:%v", coin.amount)
}
func (coin *Coin) Destroy() {
	setEmptyObject(coin.Id())
}

func NewCoin(amount int, id int) *Coin {
	return &Coin{amount: amount, id: id}
}
