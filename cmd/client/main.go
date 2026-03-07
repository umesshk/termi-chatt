package main

import (
	"github.com/gorilla/websocket"
	"fmt"
	"os"
	"log"
	"bufio"
	"strings"
	"encoding/json"
)


 type UserMessage struct {
 Msgtype 		string  `json:"type"` 
 Message		string  `json:"message,omitempty"` 
 RoomId   	int			`json:"roomId,omitempty"`
 Username 	string 	`json:"username"`
}

type ServerResponse struct {
	Type   			string 	`json:"type"`
	UserName 		string 	`json:"username"`
	Message     string  `json:"message"`
	RoomId			int 		`json:"roomId"`
}


func CreateConnection() (*websocket.Conn ,  error)  {
	
		conn , _ ,  err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)

		if err != nil {
			return nil, err  
		}
		
		log.Println(" Connected to Server... ")
		fmt.Printf("\n")

   return conn,nil 

		

}


func GetUserInput(userInput chan string,done chan struct{} ){
		reader := bufio.NewReader(os.Stdin)


		for {

			select {

				case <- done : 
					return 
		
					default: 
						msg , _ := reader.ReadString('\n')
						msg = strings.TrimSpace(msg)
	
						if msg != ""{
							userInput <- msg 
					}
			}
	}

}

func GetServerResponse(conn *websocket.Conn, serverResponseChan chan ServerResponse, done chan struct{}){

	for {

		select {

		case <- done : 
		return 
		
	default: 

		_, res, err := conn.ReadMessage()
	
		if err != nil {
		return 
		}
		
		var server_response ServerResponse
		
		if err := json.Unmarshal(res,&server_response); err != nil {
			return
		}
		

	serverResponseChan <- server_response
}
}
}








func StartConnection(Type string){

	var user_name string
	fmt.Printf("Enter your name :  ")
	fmt.Scanln(&user_name)

		var room_id int

	if Type=="join" {
		
		fmt.Printf("Enter the Room Id : ")
		fmt.Scanln(&room_id)

	}

	conn, err := CreateConnection()

	if err != nil {
		fmt.Println(err)
		return 
	}
	
 if Type == "create"{
	 if err := conn.WriteJSON(UserMessage{Msgtype:Type,Username:user_name}); err != nil {
		fmt.Println(err)
	  return 	
	 }
 }

 if Type == "join"{

	 if err := conn.WriteJSON(UserMessage{Msgtype:Type, RoomId:room_id, Username: user_name}); err != nil {
			fmt.Println(err)
			return 
	 }
 }

  userInput := make(chan string)
	serverResponseChan := make(chan ServerResponse)
	done := make(chan struct{})

	go GetUserInput(userInput,done)
	go GetServerResponse(conn,serverResponseChan,done)

for {
	select {
	 	
		case msg :=  <- userInput: 

    fmt.Print("\033[A\033[K")
		
		if msg == "/leave"{
			fmt.Println("Exiting....")
			if err := conn.WriteJSON(UserMessage{Msgtype:"leave",Username:user_name,RoomId:room_id}); err != nil {
		}
		  close(done)
			conn.Close()
			return 
	
		}
		conn.WriteJSON(UserMessage{
				Msgtype: "message",
				Username: user_name,
				RoomId: room_id,
				Message: msg,
			})
		case 	msg :=<- serverResponseChan: 

			switch msg.Type  {
			
					case "room_created" :
					
						  room_id = msg.RoomId	
							fmt.Printf("\r\033[K%s\n",msg.Message)

					case "room_joined": 

							fmt.Printf("\r\033[K%s\n",msg.Message)
				  
					case "error" : 

							fmt.Printf("\r\033[K%s\n",msg.Message)

					
					case "chat_message": 	
							fmt.Printf("\r\033[K%s\n",msg.Message)
			
					case "leave": 
							fmt.Printf("\r\033[K%s\n",msg.Message)
					 
					default: 
					 		fmt.Println("No Valid Response From Server")
					}
							
							fmt.Println( "----------------------- ")
							fmt.Printf("Enter Message : ")
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
		StartConnection("create")	
 	
	case 2: 
		fmt.Println("Joining Room ")
		StartConnection("join")	
	 
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
