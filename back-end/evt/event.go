package evt

import (
	"../gs"
	"../obj"
	"../typ"
	"time"
)

const EVENT = "event"

type Event struct {
	typ.Type

	Emitter  obj.Objectable
	Receiver obj.Objectable

	FromCoord gs.Coord
	ToCoord   gs.Coord

	EventType string
	Timestamp time.Time
}

func NewEvent(emitter obj.Objectable, receiver obj.Objectable, from gs.Coord, to gs.Coord, eventType string) *Event {
	return &Event{
		Emitter:   emitter,
		Receiver:  receiver,
		FromCoord: from,
		ToCoord:   to,
		EventType: eventType,
		Timestamp: time.Now(),
		Type:      typ.NewType(EVENT),
	}
}

func LoadEvent() typ.Typical {
	return &Event{Type: typ.NewType(EVENT)}
}
