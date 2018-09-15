package wo

import (
	"../gs"
	"../items"
	"../terr"
	"math/rand"
)

func InitializeWorld() {
	for i := 0; i < gs.WORLD_SIZE; i++ {
		for j := 0; j < gs.WORLD_SIZE; j++ {
			coord := gs.NewCoord(i, j)
			coinChance := rand.Intn(10)
			rockChance := rand.Intn(8)

			if coinChance == 0 {
				coins := items.NewCoinStack()
				AddItemsToStack(coord, coins)
			} else if rockChance == 0 {
				SetTerrain(coord, terr.NewRock())
			}
		}
	}
}
