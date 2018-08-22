package gs

type Vector struct {
	X int
	Y int
}

func NewVector(x int, y int) Vector {
	return Vector{X: x, Y: y}
}
