package server

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,

		CheckOrigin: func(r *http.Request) bool {return true},

}


func Handlehttp(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello from the server ")
}

func Ws_Handler(w http.ResponseWriter, r* http.Request){
	 conn, err := upgrader.Upgrade(w,r,nil)

	 if err!=nil {
		log.Println("Error Upgrading to Websockets ")
	 }

	defer conn.Close()

	for {
		messageType, p, err := conn.ReadMessage()

		if err != nil {
			log.Fatal("Error Reading message from %s", r.RemoteAddr)
		}

		log.Println("Message Recived from the client : ")
		fmt.Println(string(p))

		if err := conn.WriteMessage(messageType, p ); err!=nil {
			log.Fatal("Error Writing message")
		}else {
			log.Println("Message Writeen Succesfully ")
		}		
	}
}

func StartServer() {
	http.HandleFunc("/", Handlehttp)
	http.HandleFunc("/ws",Ws_Handler)
	log.Println("Server is running at port :8080")
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Println("An Error Occured starting server")
	}
}
