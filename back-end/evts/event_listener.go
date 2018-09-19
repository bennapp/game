package evts

import (
	"../evt"
	"../obj"
	"../rc"
	"../typf"
)

func EventListener(object obj.Objectable) chan *evt.Event {
	eventChannel := make(chan *evt.Event)
	serializedEventChanel := rc.Manager().SubscribeToObjectChannel(object)

	go func() {
		for {
			serializedEvent := <-serializedEventChanel
			eventChannel <- typf.Factory().DeserializeEvent(serializedEvent)
		}
	}()

	return eventChannel
}
