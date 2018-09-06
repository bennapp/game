package el

const ROCK = "rock"

type Rock struct {
	Type string
}

func (rock *Rock) String() string {
	return "R"
}
