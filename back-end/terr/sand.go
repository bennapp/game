package terr

const SAND = "sand"

func NewSand() Terrain {
	permeable := true
	return newTerrain(permeable, SAND, 0.75)
}
