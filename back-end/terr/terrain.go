package terr

import (
	"../typ"
)

type Terrain struct {
	typ.Type
	Permeable bool
	Friction  int
}

type TerrainElement interface {
	Terrainable() bool
}

// Stub to force all terrains to be a TerrainElement interface
func (terrain *Terrain) Terrainable() bool {
	return true
}

func newTerrain(terrainType string, permeable bool) Terrain {
	return Terrain{Type: typ.NewType(terrainType), Permeable: permeable}
}

func loadTerrain(terrainType string) Terrain {
	return Terrain{Type: typ.NewType(terrainType)}
}
