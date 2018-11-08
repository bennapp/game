package evt

import (
	"../gs"
	"../typ"
	"time"
)

const EVENT = "event"

type Event struct {
	typ.Type

	EmitterId  string
	ReceiverId string

	FromCoord gs.Coord
	ToCoord   gs.Coord

	EventType string
	Timestamp int64
}

func NewEvent(emitterId string, receiverId string, from gs.Coord, to gs.Coord, eventType string) *Event {
	return &Event{
		EmitterId:  emitterId,
		ReceiverId: receiverId,
		FromCoord:  from,
		ToCoord:    to,
		EventType:  eventType,
		Timestamp:  time.Now().Unix(),
		Type:       typ.NewType(EVENT),
	}
}

func LoadEvent() typ.Typical {
	return &Event{}
}
