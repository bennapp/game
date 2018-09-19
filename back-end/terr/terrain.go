package terr

type Terrain struct {
	TerrainType string
	Permeable   bool
	Friction    float64
}

func newTerrain(permeable bool, terrainType string, friction float64) Terrain {
	return Terrain{Permeable: permeable, TerrainType: terrainType, Friction: friction}
}
