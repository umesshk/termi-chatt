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
 Username 	string 	`json:"username"`
}

func CreateConnection() (*websocket.Conn ,  error)  {
	
		conn , _ ,  err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)

		if err != nil {
			return nil, err  
		}
		
		log.Println(" Connected to Server ")

   return conn,nil 

		

}


 var serverResponse = make(chan bool )
func GetServerMessage (conn *websocket.Conn){
				
				for {
			
					_,p,err := conn.ReadMessage()

					if err != nil {
						log.Println(err)
						continue
					}
			
					fmt.Printf("\n%s\n",string(p))

					serverResponse <- true

				}
}



func CreateRoom(){
  user_name := "Makito"

	conn , err := CreateConnection() 

    if err != nil {
			fmt.Println("Error Creating Connection ",err)
		}

			go GetServerMessage(conn) 

		if err := conn.WriteJSON(UserMessage{Msgtype: "create",Username:user_name}); 
 err != nil {
			log.Println("Error Occured", err)
			
			}
			
		
 		 <-serverResponse	

			for {
			
			var user_message string 
			fmt.Printf("Enter Message : ")
			fmt.Scanln(&user_message)

			if user_message != ""{
				
			 if err := conn.WriteJSON(UserMessage{Msgtype:"message",Username:user_name,Message:user_message,RoomId:1}); err != nil {
					log.Println(err)
					continue
				}

			}

		}

}

func JoinRoom(){

	var user_name string  
	fmt.Printf("Enter Your User Name : ")
	fmt.Scanln(&user_name)
	
	var roomId int
	fmt.Printf("Enter Room Id: ")
	fmt.Scanln(&roomId)


 conn, err := CreateConnection()

 			if err != nil {
			fmt.Println(err)
			return 
			}
 
			go GetServerMessage(conn) 

			if err := conn.WriteJSON(UserMessage{Msgtype: "join",RoomId: 1 , Username:user_name}); 
 err != nil {
			log.Println("Error Occured", err)
			
			}
			
		
 		 <-serverResponse	

			for {
			
			var user_message string 
			fmt.Printf("Enter Message : ")
			fmt.Scanln(&user_message)

			if user_message != ""{
				
			 if err := conn.WriteJSON(UserMessage{Msgtype:"message",Username:user_name,Message:user_message,RoomId:1}); err != nil {
					log.Println(err)
					continue
				}

			}

		}

 

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
 		JoinRoom()
	 
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
