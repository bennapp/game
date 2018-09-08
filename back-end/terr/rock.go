package terr

const ROCK = "rock"

func NewRock() Terrain {
	permeable := false
	return newTerrain(permeable, ROCK)
}
