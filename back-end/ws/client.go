package ws

import (
	"../cell"
	"../dbs"
	"../evts"
	"../gs"
	"../movs"
	"../obj"
	"../wo"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/vmihailenco/msgpack"
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

	closed chan bool
}

type playerEvent struct {
	Type string
	From gs.Coord
	To   gs.Coord
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump(player *obj.Player) {
	defer func() {
		c.closed <- true
		c.hub.unregister <- c
		c.conn.Close()
	}()
	// c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.conn.ReadMessage()

		if err != nil {
			log.Println(err)
			break
		}

		event := &playerEvent{}
		err = msgpack.Unmarshal(message, event)

		fmt.Println(event)

		if err != nil {
			log.Println(err)
			break
		}

		log.Println(fmt.Sprintf("Move Event %v,%v => %v,%v", event.From.X, event.From.Y, event.To.X, event.To.Y))

		moveVector := gs.NewVector(0, 0)

		moveVector.X = event.To.X - player.GetLocation().X
		moveVector.Y = event.To.Y - player.GetLocation().Y

		moveVector = moveVector.Normalize()

		fmt.Printf("move vector in client: %v \n", moveVector)

		movs.RegulateMove(player, moveVector)

		// The is where we could correct the client
		// we would need to detect if the client is saying it is moving from an invalid location
		// then we would send correction like below

		//gameState := gameStateMapping{}
		//
		//gameState["globalPlayerLocation"] = player.GridCoord
		//
		//gameStateAsString, _ := json.Marshal(gameState)
		//c.send <- []byte(gameStateAsString)
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
		c.closed <- true
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:

			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.BinaryMessage)
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

type dboLookup map[string]*cell.Cell
type gameStateMapping map[string]interface{}

func (c *Client) beamState(player *obj.Player) {
	gameState := gameStateMapping{}
	coordinateMapping := dboLookup{}

	halfWidth := gs.LOADED_VISION_DISTANCE / 2
	v := gs.NewVector(-halfWidth, -halfWidth)

	for i := 0; i < gs.LOADED_VISION_DISTANCE; i++ {
		for j := 0; j < gs.LOADED_VISION_DISTANCE; j++ {
			coord := player.GetLocation().AddVector(v)
			cell := dbs.LoadCell(coord)

			if !cell.IsEmpty() {
				coordinateMapping[coord.Key()] = cell
			}
			v.X += 1
		}
		v.X = -halfWidth
		v.Y += 1
	}

	gameState["coordinates"] = coordinateMapping
	gameStateAsString, _ := msgpack.Marshal(gameState)

	c.send <- []byte(gameStateAsString)
}

func (c *Client) beamInitialState(player *obj.Player) bool {
	select {
	case <-c.closed:
		return true
	default:
		gameState := gameStateMapping{}

		gameState["globalPlayerLocation"] = player.GetLocation()
		gameState["playerId"] = player.ObjectId()

		gameStateAsString, _ := msgpack.Marshal(gameState)
		c.send <- []byte(gameStateAsString)
	}

	return false
}

func (c *Client) beamStateUntilClosed(player *obj.Player) {
	closed := c.beamInitialState(player)

	if closed {
		return
	}

	c.beamState(player)

	// beam state should send location where abouts the first time,
	// from then on it should not send player locations

	//for {
	//	select {
	//	case <-c.closed:
	//		return
	//	default:
	//		c.beamState(player)
	//		time.Sleep(500 * time.Millisecond)
	//	}
	//}
}

func (c *Client) relayEvents(player *obj.Player) {
	eventChannel := evts.EventListener(player)

	for {
		event := <-eventChannel
		eventAsString, _ := msgpack.Marshal(event)
		c.send <- []byte(eventAsString)
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

	player := wo.CreatePlayer()
	onClose := func(code int, text string) error {
		wo.DeletePlayer(player)
		return nil
	}
	client.conn.SetCloseHandler(onClose)

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump(player)
	go client.beamStateUntilClosed(player)
	go client.relayEvents(player)
}
