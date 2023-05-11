package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrager = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

func homePage(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "Home page")
}

func reader (connection *websocket.Conn) {
	for {
		messageType, message, err := connection.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		log.Println(string(message))

		if err := connection.WriteMessage(messageType, message); err != nil {
			log.Println(err)
			return
		}
	}
}

func wsEndpoint(writer http.ResponseWriter, request *http.Request) {
	upgrager.CheckOrigin = func(request *http.Request) bool {return true}

	ws, err := upgrager.Upgrade(writer, request, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client has been connected")

	err = ws.WriteMessage(1, []byte("Hi, client"))
	if err != nil {
		log.Println(err)
	}

	reader(ws)
}

func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndpoint)
}

func main() {
	// todo: log this
	log.Println("WS server has been started")
	setupRoutes()
	log.Println(http.ListenAndServe(":4000", nil))
}
