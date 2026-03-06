package main

import (
	"github.com/gorilla/websocket"
	"log"
	"encoding/json"
	"net/http"
	"fmt"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}


type UserMessage struct {
 Msgtype 		string  `json:"type"` 
 Message		string  `json:"message,omitempty"` 
 RoomId   	int			`json:"roomId,omitempty"`
 Username 		string 					`json:"username"`
}

type User struct {
	UserId   		int 						`json:"userId"`
	Username 		string 					`json:"username"`
  User_conn 				*websocket.Conn `json:"conn,omitempty"`
}

var roomId int = 0
var userId int = 0 
//  roomid -> websocket connection slice
var room_map = make(map[int][]User)
// websocket connection -> room id 
var conn_map = make(map[*websocket.Conn]int)


func MainHanlder(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	defer conn.Close()

	if err != nil {
		log.Fatal(err)
	 return 	
	}

	for {
		messageType, p, err := conn.ReadMessage()

		if err != nil {
			log.Println(err)
			continue
		}

	var ClientMessage UserMessage 
 
	if err := json.Unmarshal(p,&ClientMessage); err != nil{
			log.Println("Error During Parsing " , err)
			
			continue	
	} 
		switch ClientMessage.Msgtype {
		case "create" :
				roomId++
				userId++
				user_name := ClientMessage.Username
				user := User{userId,user_name,conn}

				room_map[roomId] = append(room_map[roomId], user)	
				userId++

	    	conn_map[conn] = roomId	

				message := fmt.Sprintf("User %v Created Room %v",user_name,roomId)

				if err := conn.WriteMessage(messageType,[]byte(message)); err!=nil {
					log.Println(err)
					continue 
				}	 

				log.Printf("User  Username : %v  Created room %v\n" ,  user_name, roomId)
	
			case "join" :  
			
				room_id   := ClientMessage.RoomId
				user_name := ClientMessage.Username
				userId++ 
				user := User{userId, user_name, conn}

		   
				log.Println("Client Room Id ", room_id)
				_, ok := room_map[room_id]	
					if !ok {
						if err := conn.WriteMessage(messageType,[]byte("Room doesn't exist")); err != nil {
							fmt.Println(err)
							continue
						}

					}else {

					   	conn_room_id , ok := conn_map[conn]

							if ok && conn_room_id == room_id {
								if err := conn.WriteMessage(messageType,[]byte("Already in  the room... ")); err != nil {
									log.Printf("User Already in Room ")
								} 

						}else {
						room_map[room_id] = append(room_map[room_id], user)
	    			conn_map[conn] = room_id
						message := fmt.Sprintf("%v  Joined room : %v", user_name,room_id)

						room_conn := room_map[room_id]
					

					for _,receiver:= range room_conn {
							 receiver_name := receiver.Username
						 	 receiver_conn := receiver.User_conn
					 			
					 		if err := receiver_conn.WriteMessage(messageType, []byte(message)); err != nil {
										log.Println("Error Sending Message ")
							continue
					 } 
						 log.Printf("Message Written to User's :  %v \n",receiver_name)
				 
				 
						if err := conn.WriteMessage(messageType,[]byte(message)); err != nil {
						log.Println(err)
				  	continue 	
						}
					log.Printf(" %v Joined Room %v\n",user_name, room_id)
					}
					}
				}
	    
		case "message":
			 
				room_id := ClientMessage.RoomId
			  sender_name := ClientMessage.Username	
				if room_id == 0 {
					log.Println("No RoomId Provided ")
					if err := conn.WriteMessage(messageType,[]byte("Room Id is Required  ")) ; err != nil {
					log.Println(err)
					continue
				}
			}

				

				join_room_id , ok := conn_map[conn]

				if !ok {
					fmt.Println("Connection not Found ")
					if err := conn.WriteMessage(messageType,[]byte("Please Join the room first ")) ; err != nil {
					log.Println(err)
					continue
				}

			}else {
				if join_room_id != room_id {
					fmt.Println("Wrong Room")
						if err := conn.WriteMessage(messageType,[]byte("Please Join the room first ")) ; err != nil {
						log.Println(err)
						continue
						}	
				}
 				 current_room := room_map[join_room_id]
			 	 sender_message := ClientMessage.Message

				 log.Printf("Sender Message %v\n", sender_message)

									
				 for _,receiver:= range current_room {
					 receiver_name := receiver.Username
					 receiver_conn := receiver.User_conn
					 fmt.Println("Writing Message to all user ")
				 		message_to_send := fmt.Sprintf(" %v : %v ",sender_name,sender_message)
					 if err := receiver_conn.WriteMessage(messageType, []byte(message_to_send)); err != nil {
						log.Println("Error Sending Message ",err)
						continue
				 }
				 log.Printf("Message Written to User's :  %v \n",receiver_name)
				 


				 }

			}

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
