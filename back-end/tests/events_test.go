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
	playerId := player.ObjectId()

	dbs.SaveObject(player)
	dbs.SaveObjectLocation(player.GetLocation(), player)
	eventChannel := evts.EventListener(player)

	go func() {
		for {
			event := <-eventChannel
			t.Log(event)
		}
	}()

	coord := player.GetLocation()

	event := evt.NewEvent(playerId, playerId, coord, coord, "wave hi")
	evts.EmitEvent(event)

	newCoord := coord.AddVector(gs.Vector{X: 4, Y: 4})
	otherCoord := coord.AddVector(gs.Vector{X: 2, Y: 2})
	eventBye := evt.NewEvent(playerId, playerId, newCoord, otherCoord, "wave bye")
	evts.EmitEvent(eventBye)
}
