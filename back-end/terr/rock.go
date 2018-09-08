package terr

import "../typ"

const ROCK = "rock"

type Rock struct {
	Terrain
}

func NewRock() *Rock {
	permeable := false
	return &Rock{Terrain: newTerrain(ROCK, permeable)}
}

func LoadRock() typ.Typical {
	return &Rock{Terrain: loadTerrain(ROCK)}
}
