package pnt

import (
	"../terr"
)

type Paint struct {
	terr.Terrain
	CellState
}

// pvp, pve, safe zones, danger zones etc.
type CellState struct {
}
