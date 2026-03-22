package main

import (
	"github.com/gorilla/websocket"
	"log"
	"encoding/json"
	"net/http"
	"github.com/umeshhk/termi-chatt/internal/service/ws"
	"github.com/umeshhk/termi-chatt/internal/user"
	"github.com/umeshhk/termi-chatt/internal/database"
	"database/sql"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}




var db *sql.DB


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
						ws.HandleCreate(ClientMessage,conn,db)
	
		case "join" :  
						ws.HandleJoin(ClientMessage, conn,db)
	    
		case "message":
			 		  ws.HandleMessage(ClientMessage,conn,db)
		
		case "leave": 
						ws.HandleLeave(ClientMessage,conn,db)	

		default : 
				if err := conn.WriteMessage(messageType,[]byte("Invalid Input ")) ; err != nil {
					log.Println(err)
				}
			

		}

	}
}


func main() {

	PORT := ":8080"

	log.Printf("Starting Server on PORT  %v\n", PORT)
	
	var db_err error

	db,db_err = database.ConnectDatabse()

	if db_err != nil {
		log.Fatal("Database not Connected...",db_err)
	}

	if err := db.Ping(); err!=nil {
		log.Fatal("DB not reachable",err)	
	}

	log.Println("Connected to Database Succesfully")

	http.HandleFunc("/ws", MainHanlder)

	err := http.ListenAndServe(PORT, nil)

	if err != nil {
		log.Fatal(err)
	}

}
