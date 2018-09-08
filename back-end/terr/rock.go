package terr

const ROCK = "rock"

type Rock struct {
	Terrain
}

func NewRock() *Rock {
	permeable := false
	return &Rock{Terrain: newTerrain(ROCK, permeable)}
}

func loadRock() TerrainElement {
	return &Rock{Terrain: loadTerrain(ROCK)}
}
