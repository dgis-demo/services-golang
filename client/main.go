package main

import (
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

type Point struct {
	Lat       float32 `json:"lat"`
	Lon       float32 `json:"lon"`
	Magnitude float32 `json:"magnitude"`
}

func wsEndpoint(writer http.ResponseWriter, request *http.Request) {
	upgrager.CheckOrigin = func(request *http.Request) bool { return true }

	ws, err := upgrager.Upgrade(writer, request, nil)
	if err != nil {
		log.Printf("ws upgrader error: %s", err)
	}

	for {
		point := &Point{
			Lat:       float32((rand.Intn(90) - rand.Intn(90))) + rand.Float32(),
			Lon:       float32((rand.Intn(180) - rand.Intn(180))) + rand.Float32(),
			Magnitude: float32(rand.Intn(10)) + rand.Float32(),
		}

		err = ws.WriteJSON(point)
		if err != nil {
			log.Printf("JSON messaging error: %s", err)
			break
		}

		time.Sleep(time.Second / RPS)
	}
}

func setupRoutes() {
	http.HandleFunc("/", wsEndpoint)
}

func main() {
	log.Println("WS server has been started")
	setupRoutes()
	log.Println(http.ListenAndServe(":4000", nil))
}
