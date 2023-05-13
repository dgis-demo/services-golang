package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrager = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func homePage(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "Home page")
}

type Point struct {
	Lat       float32 `json:"lat"`
	Lon       float32 `json:"lon"`
	Magnitude float32 `json:"magnitude"`
}

func wsEndpoint(writer http.ResponseWriter, request *http.Request) {
	upgrager.CheckOrigin = func(request *http.Request) bool { return true }

	ws, err := upgrager.Upgrade(writer, request, nil)
	if err != nil {
		log.Println(err)
	}

	for {
		point := &Point{
			Lat:       float32((rand.Intn(90) - rand.Intn(90))) + rand.Float32(),
			Lon:       float32((rand.Intn(180) - rand.Intn(180))) + rand.Float32(),
			Magnitude: float32(rand.Intn(10)) + rand.Float32(),
		}
		marshalledPoint, err := json.Marshal(point)
		if err != nil {
			log.Println(err)
		}

		err = ws.WriteJSON(string(marshalledPoint))
		if err != nil {
			log.Println(err)
		}

		time.Sleep(time.Second / RPS)
	}
}

func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndpoint)
}

func main() {
	log.Println("WS server has been started")
	setupRoutes()
	log.Println(http.ListenAndServe(":4000", nil))
}
