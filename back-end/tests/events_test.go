package tests

import (
	"../dbs"
	"../evt"
	"../evts"
	"../gs"
	"../obj"
	"testing"
)

func TestEvents(t *testing.T) {
	player := obj.NewPlayer()
	dbs.SaveObject(player)
	dbs.SaveObjectLocation(player.Location, player)
	eventChannel := evts.EventListener(player)

	go func() {
		for {
			select {
			case event := <-eventChannel:
				t.Log(event)
			default:
				// no op
			}
		}
	}()

	coord := player.Location

	event := evt.NewEvent(player, player, coord, coord, "wave hi")
	evts.EmitEvent(event)

	newCoord := coord.AddVector(gs.Vector{X: 4, Y: 4})
	otherCoord := coord.AddVector(gs.Vector{X: 2, Y: 2})
	eventBye := evt.NewEvent(player, player, newCoord, otherCoord, "wave bye")
	evts.EmitEvent(eventBye)
}
