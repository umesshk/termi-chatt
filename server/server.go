package server

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	defer conn.Close()

	if err != nil {
		log.Fatal(err)
		return
	}

	for {
		messageType, p, err := conn.ReadMessage()

		if err != nil {
			log.Fatal(err)
			return
		}

		ClientMessage := string(p)

		log.Printf("%v", ClientMessage)

		if ClientMessage == "pong" {

			if err := conn.WriteMessage(messageType, []byte("ping")); err != nil {
				log.Println(err)
				return
			}

		} else {
			if err := conn.WriteMessage(messageType, []byte("pong")); err != nil {
				log.Println(err)
				return
			}

		}

	}

}

func StartServer() {

	PORT := ":8080"

	log.Printf("Starting Server on  %v\n", PORT)

	http.HandleFunc("/ws", HomeHandler)
	err := http.ListenAndServe(PORT, nil)

	if err != nil {
		log.Fatal(err)
	}

}
