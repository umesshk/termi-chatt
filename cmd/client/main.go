package main

import (
	"github.com/gorilla/websocket"
	"fmt"
	"os"
	"log"
	"bufio"
	"strings"
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


func GetUserInput(userInput chan string ){
		reader := bufio.NewReader(os.Stdin)


		for {
		msg , _ := reader.ReadString('\n')
		msg = strings.TrimSpace(msg)

		if msg != ""{
			userInput <- msg 
	}

		}

}

func GetServerResponse(conn *websocket.Conn, serverResponse chan string ){

	for {
	
		_, p, err := conn.ReadMessage()
	if err != nil {
		fmt.Println(err)
		continue
	}
	serverResponse <- string(p)
		
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
	serverResponse := make(chan string)

	go GetUserInput(userInput)
	go GetServerResponse(conn,serverResponse)

for {
	select {
	 	
		case msg :=  <- userInput: 

    fmt.Print("\033[A\033[K")
		
		if msg == "/exit"{
			fmt.Println("Exiting....")
			os.Exit(0)
		}
		conn.WriteJSON(UserMessage{
				Msgtype: "message",
				Username: user_name,
				RoomId: room_id,
				Message: msg,
			})
		case 	msg :=<- serverResponse: 
				fmt.Printf("\r\033[k%s\n",msg)
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
