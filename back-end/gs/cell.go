package gs

import "sync"

type Cell struct {
	mux   sync.Mutex
	coord Coord
}

func (cell *Cell) Mux() *sync.Mutex {
	return &cell.mux
}

func (cell *Cell) Coord() Coord {
	return cell.coord
}
