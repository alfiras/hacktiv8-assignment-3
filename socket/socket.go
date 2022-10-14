package socket

import (
	"fmt"
	"hacktiv8-assignment-3/data"

	"github.com/gorilla/websocket"
)

type SocketPayload struct {
	Message string
}

var AllConnections = make([]*websocket.Conn, 0)

func AddConnection(conn *websocket.Conn) {
	AllConnections = append(AllConnections, conn)
	fmt.Printf("Connection added -> length now: %d\n", len(AllConnections))
}

func RemoveConnection(conn *websocket.Conn) {
	var ind int
	for i, c := range AllConnections {
		if c == conn {
			ind = i
			break
		}
	}
	AllConnections = append(AllConnections[:ind], AllConnections[ind+1:]...)
	fmt.Printf("Connection removed -> length now: %d\n", len(AllConnections))
}

func Inject() {
	data.AddDep(func(d data.Data) {
		for _, c := range AllConnections {
			c.WriteJSON(d)
		}
		fmt.Printf("data %+v\n", d)
	})
}
