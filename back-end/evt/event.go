package evt

import (
	"../gs"
	"../obj"
	"time"
)

type Event struct {
	Emitter  obj.Objectable
	Receiver obj.Objectable

	FromCoord gs.Coord
	ToCoord   gs.Coord

	EventType string
	Timestamp time.Time
}

func NewEvent(emitter obj.Objectable, receiver obj.Objectable, from gs.Coord, to gs.Coord) *Event {
	return &Event{Emitter: emitter, Receiver: receiver, FromCoord: from, ToCoord: to, Timestamp: time.Now()}
}
