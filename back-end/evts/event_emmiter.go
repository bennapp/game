package evts

import (
	"../dbs"
	"../evt"
	"../gs"
	"../obj"
	"../rc"
)

func EmitEvent(event *evt.Event) {
	objectLocationCordMappings := make(map[gs.Coord]bool)

	halfVisionDistance := gs.LOADED_VISION_DISTANCE / 2
	v := gs.NewVector(-halfVisionDistance, -halfVisionDistance)

	for i := 0; i < gs.LOADED_VISION_DISTANCE; i++ {
		for j := 0; j < gs.LOADED_VISION_DISTANCE; j++ {
			objectLocationCordMappings[event.FromCoord.AddVector(v)] = true
			objectLocationCordMappings[event.ToCoord.AddVector(v)] = true

			v.X += 1
		}
		v.X = -halfVisionDistance
		v.Y += 1
	}

	for coord, _ := range objectLocationCordMappings {
		object := dbs.LoadObjectByCoord(coord)

		switch player := object.(type) {
		case *obj.Player:
			rc.Manager().WriteToObjectEventChannel(player, event)
		}
	}
}
