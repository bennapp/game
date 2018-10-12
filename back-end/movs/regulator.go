package movs

import (
	"../dbs"
	"../gs"
	"../mov"
	"time"
)

func RegulateMove(movable mov.Movable, vector gs.Vector) {
	if movable.GetMovesToRegulate() == nil {
		ch := make(chan gs.Vector)
		movable.SetMovesToRegulate(ch)

		go relayMovesToBuffer(movable)
		go readMovesAtRegulatedInterval(movable)
	}

	enqueueMoves(movable, vector)
}

func enqueueMoves(movable mov.Movable, vector gs.Vector) {
	movable.GetMovesToRegulate() <- vector
}

func relayMovesToBuffer(movable mov.Movable) {
	for {
		move := <-movable.GetMovesToRegulate()
		movable.AppendMoveBuffer(move)
	}
}

func readMovesAtRegulatedInterval(movable mov.Movable) {
	sleepOffset := time.Duration(10) * time.Millisecond

	for {
		move, pop := movable.PopMoveBuffer()

		if pop == false {
			time.Sleep(sleepOffset)
			continue
		}

		friction := 1.0
		paint := dbs.LoadPaintByCoord(movable.GetLocation())
		if paint != nil {
			friction = paint.Terrain.Friction
		}

		MoveObject(movable, move)

		cellsPerSecond := movable.GetVelocity() * friction
		secondsPerCell := 1.0 / cellsPerSecond

		milliSecondsPerCell := secondsPerCell * 1000
		sleepTime := time.Duration(milliSecondsPerCell) * time.Millisecond
		sleepTime -= sleepOffset

		time.Sleep(sleepTime)
	}
}
