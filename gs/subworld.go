package gs

type SubWorld struct {
	coord Coord
	grid  Grid
}

func (subWorld *SubWorld) Coord() Coord {
	return subWorld.coord
}

func (subWorld *SubWorld) Grid() Grid {
	return subWorld.grid
}

func NewSubWorld(coord Coord) SubWorld {
	return SubWorld{coord: coord}
}
