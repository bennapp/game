package gs

type World struct {
	subWorlds [WORLD_SIZE][WORLD_SIZE]SubWorld
}

func (world *World) SubWorlds() *[WORLD_SIZE][WORLD_SIZE]SubWorld {
	return &world.subWorlds
}

func NewWorld(subWorlds [WORLD_SIZE][WORLD_SIZE]SubWorld) World {
	return World{subWorlds: subWorlds}
}
