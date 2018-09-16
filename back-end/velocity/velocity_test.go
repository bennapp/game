package velocity_test

import (
	"../obj"
	"../velocity"
	"testing"
)

func TestConstants(t *testing.T) {
	playerVelocity := velocity.Constants(obj.PLAYER)

	if playerVelocity != 5.0 {
		t.Errorf("could not retreive player velocity, got value %v \n", playerVelocity)
	}
}
