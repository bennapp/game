package ws

import (
	"../wo"
		"time"
	"encoding/json"
	)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

// There is where we should beam the world
func (h *Hub) beamState() {
	wo.Init()

	id := 386
	char := "M"

	player := wo.LoadPlayer(id)
	player.Avatar = char

	//v := gs.NewVector(-5, -5)
	//visionDistance := 11

	//for i := 0; i < visionDistance; i++ {
	//	for j := 0; j < visionDistance; j++ {
	//		element, valid := wo.NextElement(player.SubWorldCoord, player.GridCoord, v)
	//		if valid {
	//			fmt.Printf("%v ", element.String())
	//		}
	//		v.X += 1
	//	}
	//	v.X = -5
	//	fmt.Println()
	//	v.Y += 1
	//}

	//	{
	//    globalPlayerLocation: {
	//      x: '2',
	//      y: '2',
	//    },
	//    coordinates: {
	//      "0,1": { type: 'coin', id: '33' },
	//      "3,4": { type: 'rock', id: '-1' },
	//      "1,1": { type: 'rock', id: '-1' },
	//    },
	//    objects: {
	//      // player: {
	//      //   "1": {
	//      //     hp: "10",
	//      //     alive: "true",
	//      //     coinCount: "22",
	//      //   },
	//      //   "2": {
	//      //     hp: "7",
	//      //     alive: "true"
	//      //   }
	//      // },
	//      coin: {
	//        "33": {
	//          amount: "11",
	//        },
	//      },
	//      rock: {
	//        "-1": {}
	//      }
	//    },
	//  };

	gameState := map[string]string{}

	globalPlayerLocation := map[string]string{}
	globalPlayerLocation["x"] = "10"
	globalPlayerLocation["y"] = "10"

	globalPlayerLocationString, _ := json.Marshal(globalPlayerLocation)
	gameState["globalPlayerLocation"] = string(globalPlayerLocationString)

	gameStateAsString, _ := json.Marshal(gameState)

	for {
		h.broadcast <- []byte(gameStateAsString)
		time.Sleep(1000 * time.Millisecond)
	}
}
