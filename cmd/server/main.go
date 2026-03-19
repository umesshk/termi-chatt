package main

import (
	"github.com/gorilla/websocket"
	"log"
	"encoding/json"
	"net/http"
	"fmt"
	"sync"	
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}



var roomId int = 0
var userId int = 0 
//  roomid -> websocket connection slice
var room_map = make(map[int][]User)
// websocket connection -> room id 
var conn_map = make(map[*websocket.Conn]int)







func MainHanlder(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Fatal(err)
	 return 	
	}
 defer conn.Close()
 
	for {
		messageType, p, err := conn.ReadMessage()

		if err != nil {
			log.Println(err)
			continue
		}

	var ClientMessage UserMessage 
 
	if err := json.Unmarshal(p,&ClientMessage); err != nil{
			log.Println("Error During Parsing " , err)
			
			continue	
	} 
		switch ClientMessage.Msgtype {

		case "create" :
						handleCreate(ClientMessage,conn)
	
		case "join" :  
						handleJoin(ClientMessage, conn)
	    
		case "message":
			 		 handleMessage(ClientMessage,conn)
		
		case "leave": 
					handleLeave(ClientMessage,conn)	

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
