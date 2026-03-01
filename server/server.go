package server

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"fmt"
	"strconv"
	"strings"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var roomId int = 0

var room_map = make(map[int]*websocket.Conn)


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

		ClientMessage := strings.TrimSpace(string(p))


		if ClientMessage == "create" {
			roomId++
			room_map[roomId] = conn	
		
			message := fmt.Sprintf("Room Created with Id : %d",roomId)

			if err := conn.WriteMessage(messageType,[]byte(message)); err!=nil {
				log.Fatal(err)
				return 
			} 

			log.Printf("User %v Created room %v " , conn.RemoteAddr(), roomId)

		} else {

			roomIdString:= ClientMessage

			if err != nil{
				log.Fatal(err)
				return 
			}
			room_id,err := strconv.Atoi(roomIdString)
			
			if err != nil{
				log.Fatal(err)
				return 
			}
			
			value , ok := room_map[room_id]	
				if !ok {
					if err := conn.WriteMessage(messageType,[]byte("Room doesn't exist")); err != nil {
						log.Fatal(err)
						return 
					}
				}else {
					fmt.Println(value)
					if err := conn.WriteMessage(messageType,[]byte("Room Joined ")); err != nil {
					log.Fatal(err)
					return 
					}
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
