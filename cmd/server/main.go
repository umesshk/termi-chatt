package main

import (
	"github.com/gorilla/websocket"
	"log"
	"encoding/json"
	"net/http"
	"fmt"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var roomId int = 0

var room_map = make(map[int][]*websocket.Conn)
var conn_map = make(map[*websocket.Conn]int)

type UserMessage struct {
 Msgtype 		string  `json:"type"` 
 Message		string  `json:"message,omitempty"` 
 RoomId   	int			`json:"roomId,omitempty"`
}

func MainHanlder(w http.ResponseWriter, r *http.Request) {

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
			continue
		}

	var ClientMessage UserMessage 
 
	if err := json.Unmarshal(p,&ClientMessage); err != nil{
			log.Println(err)
			continue	
	} 
		switch ClientMessage.Msgtype {
		case "create" :
				roomId++
				room_map[roomId] = append(room_map[roomId], conn)	
	    	conn_map[conn] = roomId	

				message := fmt.Sprintf("Room Created with Id : %d",roomId)

				if err := conn.WriteMessage(messageType,[]byte(message)); err!=nil {
					log.Println(err)
					continue 
				}	 

				log.Printf("User %v Created room %v \n" , conn.RemoteAddr(), roomId)
	
			case "join" :  
			
				room_id := ClientMessage.RoomId

				if err != nil{
					log.Fatal(err)
					return 
				}		
			
				_, ok := room_map[room_id]	
					if !ok {
						if err := conn.WriteMessage(messageType,[]byte("Room doesn't exist")); err != nil {
							log.Fatal(err)
							continue
						}

					}else {
						room_map[room_id] = append(room_map[room_id], conn)
	    			conn_map[conn] = room_id
						if err := conn.WriteMessage(messageType,[]byte("Room Joined ")); err != nil {
						log.Fatal(err)
				  	continue 	
						}
						log.Printf("User Joined Room %v\n", room_id)
					}
	    case "message": 
			
			default : 
				if err := conn.WriteMessage(messageType,[]byte("Invalid Input ")) ; err != nil {
					log.Println(err)
				}
			

		}

	}
}

func main() {

	PORT := ":8080"

	log.Printf("Starting Server on  %v\n", PORT)

	http.HandleFunc("/ws", MainHanlder)
	err := http.ListenAndServe(PORT, nil)

	if err != nil {
		log.Fatal(err)
	}

}
