package terr

const GRASS = "grass"

func NewGrass() Terrain {
	permeable := true
	return newTerrain(permeable, GRASS, 1)
}
