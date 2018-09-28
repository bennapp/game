package gs

type Vector struct {
	X int
	Y int
}

func NewVector(x int, y int) Vector {
	return Vector{X: x, Y: y}
}

func (vector *Vector) Scale(magnitude int) Vector {
	x := vector.X * magnitude
	y := vector.Y * magnitude

	return NewVector(x, y)
}
