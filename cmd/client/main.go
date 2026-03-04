package main

import (
	"github.com/gorilla/websocket"
	"fmt"
	"os"
	"log"
_	"encoding/json"
)


 type UserMessage struct {
 Msgtype 		string  `json:"type"` 
 Message		string  `json:"message,omitempty"` 
 RoomId   	int			`json:"roomId,omitempty"`
}

func CreateConnection() (*websocket.Conn ,  error)  {
	
		conn , _ ,  err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)

		if err != nil {
			return nil, err  
		}
		
		log.Println(" Connected to Server ")

   return conn,nil 

		

}


func CreateRoom(){
	fmt.Println("Room Created")

	conn , err := CreateConnection() 

    if err != nil {
			fmt.Println("Error Creating Connection ",err)
		}

	if err := conn.WriteJSON(UserMessage{Msgtype: "create"}); 
 err != nil {
			log.Println("Error Occured", err)
			
			}
			
			for {
			_ ,p,err := conn.ReadMessage()
			
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(string(p))
		}

}

func JoinRoom(roomId int ){
	


}


func main() {
	options := []string{" Create Room ", " Join Room ", " Exit "}
	
 Chances := 3 	
 
for (Chances>0){
	
	for i ,option := range options {
		fmt.Printf(" [ %d ] %s \n", i+1, option)
	}
	
	var user_choice int
	fmt.Printf(" \n Select an option : ") 
	fmt.Scanln(&user_choice)
 
	switch user_choice { 
	case 1 : 
	 	fmt.Println("Creating Room ")
		CreateRoom()
 	
	case 2: 
		fmt.Println("Joining Room ")
		fmt.Printf("Enter Room Id to join: ")
		var roomId int 
		fmt.Scanln(&roomId)
 		JoinRoom(roomId)
	 
	case 3:
		fmt.Println("Exiting Program...")
		os.Exit(0)

	default  : 
		fmt.Println("Invalid Choice... ")
		Chances--
		fmt.Printf("Remaning Tries %v\n", Chances)
		
}
}
}
