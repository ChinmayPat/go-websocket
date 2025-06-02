package room

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var rooms = make(map[string][]*websocket.Conn)

type webSocket struct {
	upgrader websocket.Upgrader
}

var webSocketHandler = webSocket{
	upgrader: websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	},
}

func CreateRoom(roomName string) error {
	if rooms[roomName] != nil {
		return errors.New("room already exist")
	}
	i := 1
	for i <= 3 {
		connectionName := fmt.Sprintf("conn%d", i)
		createWsConnection(roomName, connectionName)
		i++
	}
	return nil
}

func createWsConnection(roomName string, connectionName string) {
	wsUrl := fmt.Sprintf("/%s/%s", roomName, connectionName)
	http.HandleFunc(wsUrl, func(w http.ResponseWriter, r *http.Request) {
		wsConnection, err := webSocketHandler.upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("error %s when upgrading connection to websocket", err)
			return
		}
		log.Printf("ws connected to %s", wsUrl)
		rooms[roomName] = append(rooms[roomName], wsConnection)
		broadcastMessage(roomName, wsUrl, wsConnection)
	})
}

func broadcastMessage(roomName string, wsUrl string, currentWsConnection *websocket.Conn) {
	for {
		_, msg, err := currentWsConnection.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %s", err)
			return
		}
		log.Printf("Message from client(%s) to room(%s): %s", wsUrl, roomName, msg)
		for _, wsConnection := range rooms[roomName] {
			if wsConnection != currentWsConnection {
				wsConnection.WriteMessage(websocket.TextMessage, []byte(msg))
			}
		}
	}
}
