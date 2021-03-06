package mov

import (
	"../gs"
	"sync"
)

const MAX_MOVE_BUFFER_QUEUE_LENGTH = 3

type Movable interface {
	GetVelocity() float64

	GetLocation() gs.Coord
	SetLocation(coord gs.Coord)

	GetMovesToRegulate() chan gs.Vector
	SetMovesToRegulate(chan gs.Vector)

	PopMoveBuffer() (gs.Vector, bool)
	AppendMoveBuffer(gs.Vector)
}

type Mover struct {
	location        gs.Coord
	movesToRegulate chan gs.Vector
	moveBuffer      []gs.Vector
	moveBufferMutex sync.Mutex
}

func (mover *Mover) SetLocation(coord gs.Coord) {
	mover.location = coord
}

func (mover *Mover) GetLocation() gs.Coord {
	return mover.location
}

func (mover *Mover) GetMovesToRegulate() chan gs.Vector {
	return mover.movesToRegulate
}

func (mover *Mover) SetMovesToRegulate(ch chan gs.Vector) {
	mover.movesToRegulate = ch
}

func (mover *Mover) PopMoveBuffer() (gs.Vector, bool) {
	var vector gs.Vector

	if len(mover.moveBuffer) > 0 {
		mover.moveBufferMutex.Lock()
		vector, mover.moveBuffer = mover.moveBuffer[0], mover.moveBuffer[1:]
		mover.moveBufferMutex.Unlock()

		return vector, true
	} else {
		return vector, false
	}
}

func (mover *Mover) AppendMoveBuffer(vector gs.Vector) {
	if len(mover.moveBuffer) > MAX_MOVE_BUFFER_QUEUE_LENGTH {
		return
	}

	mover.moveBufferMutex.Lock()
	mover.moveBuffer = append(mover.moveBuffer, vector)
	mover.moveBufferMutex.Unlock()
}
