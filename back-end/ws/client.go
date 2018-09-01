package ws

import (
	"../gs"
	"../rc"
	"../wo"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		c.hub.broadcast <- message
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:

			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) beamState() {
	wo.Init()

	id := 6086

	player := wo.LoadPlayer(id)

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

	gameState := map[string]rc.Dbo{}

	for {
		v := gs.NewVector(-5, -5)
		visionDistance := 11

		for i := 0; i < visionDistance; i++ {
			for j := 0; j < visionDistance; j++ {
				element, valid := wo.NextElement(player.GridCoord, v)
				nextCoord, _ := wo.SafeMove(player.GridCoord, v)

				if valid {
					if !wo.IsEmpty(nextCoord) {
						gameState[nextCoord.Key()] = element
					}
				}
				v.X += 1
			}
			v.X = -5
			fmt.Println()
			v.Y += 1
		}

		gameStateAsString, _ := json.Marshal(gameState)
		c.send <- []byte(gameStateAsString)
		time.Sleep(1000 * time.Millisecond)
	}

}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		// TODO: FIXME needs to protect against origin forgery
		return true
	}
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
	go client.beamState()
}
