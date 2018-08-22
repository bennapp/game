package el

type Building struct {
	Element
	code string
}

func (building Building) String() string {
	return building.code
}
