package terr

type Terrain struct {
	TerrainType string
	Permeable bool
	Friction  int // how slow something should move through this
}

type TerrainElement interface {
	Terrainable() bool
}

func newTerrain(permeable bool, terrainType string) Terrain {
	return Terrain{Permeable: permeable, TerrainType: terrainType}
}
