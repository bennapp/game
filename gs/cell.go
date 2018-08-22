package gs

import "sync"

type Cell struct {
	mux           sync.Mutex
	subWorldCoord Coord
	coord         Coord
}

func (cell *Cell) Mux() *sync.Mutex {
	return &cell.mux
}

func (cell *Cell) SubWorldCoord() Coord {
	return cell.subWorldCoord
}

func (cell *Cell) Coord() Coord {
	return cell.coord
}
