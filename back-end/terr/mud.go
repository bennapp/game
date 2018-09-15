package terr

const MUD = "mud"

func NewMud() Terrain {
	permeable := true
	return newTerrain(permeable, MUD, 0.5)
}
