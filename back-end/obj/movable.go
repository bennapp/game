package obj

import "../gs"

type Movable interface {
	SetLocation(coord gs.Coord)
	GetLocation() gs.Coord
}

type Mover struct {
	location gs.Coord
}

func (mover *Mover) SetLocation(coord gs.Coord) {
	mover.location = coord
}

func (mover *Mover) GetLocation() gs.Coord {
	return mover.location
}
