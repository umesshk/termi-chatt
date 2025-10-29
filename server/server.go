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
} 

func http_Hanlder(w http.ResponseWriter, r* http.Request){
		fmt.Fprintf(w,"Hello from server")
}

func WS_Handler(w http.ResponseWriter,r *http.Request){
			conn,err := upgrader.Upgrade(w,r,nil)
	
			if err!=nil {
			log.Fatal("Error Occured Upgrading to ws"); 
	    }

			defer conn.Close()

			messageType, p, err := conn.ReadMessage()
			
			if err != nil {
			log.Fatal("Error reading message")
			return 
			}
			
			log.Printf("Message Recived from client : %s ", p )

			if err := conn.WriteMessage(messageType,p); err!=nil{
				log.Fatal("Error occured Writing Message"); 
			}
}

func StartServer(){
	http.HandleFunc("/", http_Hanlder);
	http.HandleFunc("/ws",WS_Handler); 
	log.Println("Server running on port : 8080")

	http.ListenAndServe(":8080",nil) ;
}
