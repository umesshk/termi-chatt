package main

import (
	"github.com/gorilla/websocket"
	"log"
	"encoding/json"
	"net/http"
	"github.com/umeshhk/termi-chatt/internal/service/ws"
	"github.com/umeshhk/termi-chatt/internal/user"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}






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

	var ClientMessage user.UserMessage 
 
	if err := json.Unmarshal(p,&ClientMessage); err != nil{
			log.Println("Error During Parsing " , err)
			
			continue	
	} 
		switch ClientMessage.Msgtype {

		case "create" :
						ws.HandleCreate(ClientMessage,conn)
	
		case "join" :  
						ws.HandleJoin(ClientMessage, conn)
	    
		case "message":
			 		  ws.HandleMessage(ClientMessage,conn)
		
		case "leave": 
						ws.HandleLeave(ClientMessage,conn)	

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
