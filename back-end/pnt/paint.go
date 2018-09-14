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
