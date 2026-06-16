package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/umesshk/termi-chatt/internal/config"
	"github.com/umesshk/termi-chatt/internal/database"
	"github.com/umesshk/termi-chatt/internal/redisx"
	"github.com/umesshk/termi-chatt/internal/service/ws"
	"github.com/umesshk/termi-chatt/internal/user"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
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

	client := ws.NewClient(conn)
	go client.WritePump()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			hub.RemoveClient(client)
			return
		}

		var clientMessage user.UserMessage
		if err := json.Unmarshal(p, &clientMessage); err != nil {
			log.Println("Error During Parsing ", err)
			continue
		}

		switch clientMessage.Msgtype {
		case "create":
			ws.HandleCreate(clientMessage, client, db, hub)
		case "join":
			ws.HandleJoin(clientMessage, client, db, hub)
		case "message":
			ws.HandleMessage(clientMessage, client, db, hub)
		case "leave":
			ws.HandleLeave(clientMessage, client, db, hub)
		default:
			client.Enqueue(user.ServerResponse{
				Type:    "error",
				Message: "Invalid Input",
			})
			_ = messageType
		}
	}
}

func main() {
	godotenv.Load()

	cfg := config.FromEnv()
	PORT := ":" + cfg.Port

	var db_err error

	db, db_err = database.ConnectDatabse(cfg.PostgresDSN)
	if db_err != nil {
		log.Fatal("Database not Connected...", db_err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("DB not reachable", err)
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

	log.Printf("Starting Server on PORT  %v\n", PORT)
	hub = ws.NewHub(redisClient)

	http.HandleFunc("/ws", MainHanlder)

	err = http.ListenAndServe(PORT, nil)
	if err != nil {
		log.Fatal(err)
	}
}
