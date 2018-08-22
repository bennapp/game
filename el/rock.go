package el

import "fmt"

type Rock struct {
	Element
}

func (rock *Rock) Id() string {
	return fmt.Sprintf("rock")
}
func (r Rock) String() string {
	return "R"
}

func NewRock() Rock {
	return Rock{}
}
