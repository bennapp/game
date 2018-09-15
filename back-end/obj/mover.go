package obj

import "../gs"

type Mover struct {
	location gs.Coord
}

type Movable interface {
	SetLocation(coord gs.Coord)
	GetLocation() gs.Coord
}

func (mover *Mover) SetLocation(coord gs.Coord) {
	mover.location = coord
}

func (mover *Mover) GetLocation() gs.Coord {
	return mover.location
}
