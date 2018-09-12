package evts

import (
	"../el"
	"../evt"
	"../obj"
	"../rc"
)

func EventListener(object obj.Objectable) chan *evt.Event {
	eventChannel := make(chan *evt.Event)
	serializedEventChanel := rc.Manager().SubscribeToObjectChannel(object)

	go func() {
		for {
			select {
			case serializedEvent := <-serializedEventChanel:
				eventChannel <- el.Factory().DeserializeEvent(serializedEvent)
			default:
				// no op
			}
		}
	}()

	return eventChannel
}
