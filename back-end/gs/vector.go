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

func (vector *Vector) Normalize() Vector {
	x := coerceToOne(vector.X)
	y := coerceToOne(vector.Y)

	return NewVector(x, y)
}

func coerceToOne(n int) int {
	if n > 0 {
		return 1
	}

	if n < 0 {
		return -1
	}

	return 0
}
