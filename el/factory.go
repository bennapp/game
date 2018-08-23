package el

import "../rc"

var INSTANCE ElementFactory

type ElementFactory struct {
	dboManager rc.RedisManager
}

func Factory() ElementFactory {
	if &INSTANCE == nil {
		INSTANCE = ElementFactory{dboManager: rc.Manager()}
	}

	return INSTANCE
}

//TODO - does not work
func Create() Element {
	switch v := element.(type) {
	case *el.Coin:
		storeCoinCoord(subWorld, coord, v)
	case *el.Rock:
		storeRockCoord(subWorld, coord, v)
	case *el.Player:
		storePlayerCoord(subWorld, coord, v)
		storePlayer(v)
		//default:
	}
}
