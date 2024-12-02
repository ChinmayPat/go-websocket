package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

type webSocket struct {
	upgrader   websocket.Upgrader
	connection *websocket.Conn
}

var webSocketHandler = webSocket{
	upgrader: websocket.Upgrader{},
}

func (wsh *webSocket) sendMessageToWSClient(message string) {
	wsConnection := wsh.connection
	wsConnection.WriteMessage(websocket.TextMessage, []byte(message))
	if strings.TrimSpace(string(message)) == "close" {
		log.Println("Closing connection")
		if err := wsConnection.Close(); err != nil {
			log.Printf("Error closing message %s", err)
			return
		}
		return
	}
}

func handleWS(w http.ResponseWriter, r *http.Request) {
	wsConnection, err := webSocketHandler.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error %s when upgrading connection to websocket", err)
		return
	}
	webSocketHandler.connection = wsConnection
	log.Println("ws connected")
	for {
		_, msg, err := wsConnection.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %s", err)
			return
		}
		log.Printf("Message from client: %s", msg)
		// msgAcknowlegment := "Message received: " + string(msg)
		// webSocketHandler.sendMessageToClient(msgAcknowlegment)
	}
}

func handleHTTP(w http.ResponseWriter, r *http.Request) {
	indexHtml := template.Must(template.ParseFiles("index.html"))
	if err := indexHtml.Execute(w, nil); err != nil {
		log.Println("Failed to render html: ", err)
		return
	}
}

func sendMessageToWS(w http.ResponseWriter, r *http.Request) {
	data := r.PostFormValue("message-area")
	log.Printf("Message received from http: %s", data)
	webSocketHandler.sendMessageToWSClient(data)
	// tmpl, _ := template.New("message").Parse(string(data))
	// tmpl.Execute(w, nil)
}

func main() {
	http.HandleFunc("/", handleHTTP)
	http.HandleFunc("/ws", handleWS)
	http.HandleFunc("/message", sendMessageToWS)
	log.Print("Starting server...")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
