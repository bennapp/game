package pnt

import (
	"../terr"
	"../typ"
)

const PAINT = "paint"

type Paint struct {
	typ.Type
	terr.Terrain
	Zone
}

// pvp, pve, safe zones, etc.
type Zone struct {
}

func NewPaintWithEmptyTerrain() *Paint {
	zone := Zone{}
	return &Paint{Type: typ.NewType(PAINT), Zone: zone}
}

func NewPaintByTerrain(terrain terr.Terrain) *Paint {
	zone := Zone{}
	return NewPaint(terrain, zone)
}

func NewPaint(terrain terr.Terrain, zone Zone) *Paint {
	return &Paint{Type: typ.NewType(PAINT), Terrain: terrain, Zone: zone}
}

func LoadPaint() typ.Typical {
	return &Paint{}
}

func (paint *Paint) TerrainEmpty() bool {
	return paint.Terrain.TerrainType == ""
}

func (paint *Paint) SetTerrainByType(terrainType string) {
	var terrain terr.Terrain

	switch terrainType {
	case "grass":
		terrain = terr.NewGrass()
	case "mud":
		terrain = terr.NewMud()
	case "rock":
		terrain = terr.NewRock()
	case "sand":
		terrain = terr.NewSand()
	}

	paint.Terrain = terrain
}
