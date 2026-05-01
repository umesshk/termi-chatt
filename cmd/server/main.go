package main

import (
	"github.com/gorilla/websocket"
	"log"
	"encoding/json"
	"net/http"
	"github.com/umesshk/termi-chatt/internal/service/ws"
	"github.com/umesshk/termi-chatt/internal/user"
	"github.com/umesshk/termi-chatt/internal/database"
	"database/sql"
	"github.com/umesshk/termi-chatt/internal/config"
	"github.com/umesshk/termi-chatt/internal/redisx"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Terminal client connects locally; permissive origin avoids surprising failures.
		// Tighten this if you expose the server publicly.
		return true
	},
}




var db *sql.DB
var hub *ws.Hub


func MainHanlder(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println("upgrade error:", err)
		return
	}
 defer conn.Close()
 
	for {
		messageType, p, err := conn.ReadMessage()

		if err != nil {
			log.Println(err)
			hub.RemoveConn(conn)
			return
		}

	var ClientMessage user.UserMessage 
 
	if err := json.Unmarshal(p,&ClientMessage); err != nil{
			log.Println("Error During Parsing " , err)
			
			continue	
	} 
		switch ClientMessage.Msgtype {

		case "create" :
						ws.HandleCreate(ClientMessage, conn, db, hub)
	
		case "join" :  
						ws.HandleJoin(ClientMessage, conn, db, hub)
	    
		case "message":
			 		  ws.HandleMessage(ClientMessage, conn, db, hub)
		
		case "leave": 
						ws.HandleLeave(ClientMessage, conn, db, hub)

		default : 
				if err := conn.WriteMessage(messageType,[]byte("Invalid Input ")) ; err != nil {
					log.Println(err)
				}
			

		}

	}
}


func main() {

	cfg := config.FromEnv()
	PORT := ":" + cfg.Port

	log.Printf("Starting Server on PORT  %v\n", PORT)
	
	var db_err error

	db,db_err = database.ConnectDatabse(cfg.PostgresDSN)

	if db_err != nil {
		log.Fatal("Database not Connected...",db_err)
	}
	defer db.Close()

	if err := db.Ping(); err!=nil {
		log.Fatal("DB not reachable",err)	
	}

	log.Println("Connected to Database Succesfully")

	redisClient, err := redisx.New(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		log.Fatal("Redis not reachable", err)
	}
	if redisClient != nil {
		defer redisClient.Rdb.Close()
		log.Println("Connected to Redis successfully")
	} else {
		log.Println("Redis disabled (set REDIS_ADDR to enable)")
	}

	hub = ws.NewHub(redisClient)

	http.HandleFunc("/ws", MainHanlder)

	err = http.ListenAndServe(PORT, nil)

	if err != nil {
		log.Fatal(err)
	}

}
