package main

import (
	"github.com/gorilla/websocket"
	"log"
	"encoding/json"
	"net/http"
	"fmt"
	"sync"	
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var mu sync.Mutex

type UserMessage struct {
 Msgtype 		string  `json:"type"` 
 Message		string  `json:"message,omitempty"` 
 RoomId   	int			`json:"roomId,omitempty"`
 Username 	string 					`json:"username"`
}

type User struct {
	UserId   		int 						`json:"userId"`
	Username 		string 					`json:"username"`
  User_conn  *websocket.Conn `json:"conn,omitempty"`
}

type ServerResponse struct {
	Type 			string 	`json:"type"`
	UserName 	string 	`json:"username"`	
	Message 	string 	`json:"message,omitempty"`
	RoomId  	int 	 	`json:"roomId,omitempty"`
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
				
				mu.Lock()
				room_map[roomId] = append(room_map[roomId], user)	

	    	conn_map[conn] = roomId	

				mu.Unlock()


				message := fmt.Sprintf(" Created Room with Room Id :  %v",roomId)

				server_response := ServerResponse{Type:"room_created", UserName:user_name, Message: message, RoomId: roomId}
			
				conn.WriteJSON(server_response)

				log.Printf("%v  Created room %v\n" ,  user_name, roomId)
	
			case "join" :  
			
				room_id   := ClientMessage.RoomId
				user_name := ClientMessage.Username
				userId++ 
				user := User{userId, user_name, conn}
				
		   
				log.Println("Client Room Id ", room_id)
				
				mu.Lock()
				_, ok := room_map[room_id]	
				mu.Unlock()

				
				if !ok {

					message := fmt.Sprintf("Room Doesn't Exist with room Id : %v", room_id)	
					server_response := ServerResponse{Type:"error",UserName:user_name,Message:message,RoomId: room_id}
						
					conn.WriteJSON(server_response)
					fmt.Println(message)

					}else {
							mu.Lock()
					   	conn_room_id , ok := conn_map[conn]
							mu.Unlock()

							if ok && conn_room_id == room_id {

									message := fmt.Sprintf("User %v Already  in  room Id : %v",user_name, room_id)	
						
									server_response := ServerResponse{Type:"error",UserName:user_name,Message:message,RoomId: room_id}
					
									conn.WriteJSON(server_response)
									fmt.Println(message)

						}else {

						mu.Lock()
						room_map[room_id] = append(room_map[room_id], user)
	    			conn_map[conn] = room_id
						mu.Unlock()
			

						mu.Lock()
						room_conn := room_map[room_id]
						mu.Unlock()
					

						message := fmt.Sprintf("%v  Joined the room ", user_name)
					
						server_response := ServerResponse{Type:"room_joined",UserName:user_name,Message:message,RoomId: room_id}

							for _,receiver:= range room_conn {
						 		 
								receiver_conn := receiver.User_conn
					 		
								receiver_conn.WriteJSON(server_response)


				 			log.Println(message)
						}
				}
			}

	    
		case "message":
			 
				room_id := ClientMessage.RoomId
			  sender_name := ClientMessage.Username	
				
				if room_id == 0 {
						message := fmt.Sprintf("No Room Id provided... ")	
							
						server_response := ServerResponse{Type:"error",UserName:sender_name,Message:message,RoomId: room_id}
						
						conn.WriteJSON(server_response)
					  
						log.Println(message)

						continue
			
				}

				

				mu.Lock()
				joined_room_id , ok := conn_map[conn]
				mu.Unlock()

				if !ok {
					
					message := fmt.Sprintf("Please Join the Room First : %v ",room_id)	
							
						server_response := ServerResponse{Type:"error",UserName:sender_name,Message:message,RoomId: room_id}
						
						conn.WriteJSON(server_response)
					  
						log.Println(message)

						continue

			}else {

				if joined_room_id != room_id {
					message := fmt.Sprintf("Wrong Room Id Provided : %v ",room_id)	
							
					server_response := ServerResponse{Type:"error",UserName:sender_name,Message:message,RoomId: room_id}
						
						conn.WriteJSON(server_response)
					  
						log.Println(message)

						continue
			}

				mu.Lock()
 				 current_room := room_map[joined_room_id]
				 mu.Unlock()
			 	 sender_message := ClientMessage.Message

				 log.Printf("Sender Message %v\n", sender_message)

									
					fmt.Println("Writing Message to all user ")
			
					for _,receiver:= range current_room {

					 receiver_name := receiver.Username
					 receiver_conn := receiver.User_conn
					 
				 	 message_to_send := fmt.Sprintf("%v",sender_message)
					
					 server_response := ServerResponse{Type:"chat_message", UserName:sender_name,Message:message_to_send,RoomId:room_id }
						
					 receiver_conn.WriteJSON(server_response)
					 
				 log.Printf("Message Written to User :  %v \n",receiver_name)
				 


				 }

			}
		case "leave": 
				room_id := ClientMessage.RoomId
			  sender_name := ClientMessage.Username	
			
				
				if room_id == 0 {
						message := fmt.Sprintf("No Room Id provided... ")	
							
						server_response := ServerResponse{Type:"error",UserName:sender_name,Message:message,RoomId: room_id}
						
						conn.WriteJSON(server_response)
					  
						log.Println(message)

						continue
			
				}

				
				mu.Lock()
				joined_room_id , ok := conn_map[conn]
				mu.Unlock()

				if !ok {
					
					message := fmt.Sprintf("Please Join the Room First : %v ",room_id)	
							
						server_response := ServerResponse{Type:"error",UserName:sender_name,Message:message,RoomId: room_id}
						
						conn.WriteJSON(server_response)
					  
						log.Println(message)

						continue

			}else {

				if joined_room_id != room_id {
					message := fmt.Sprintf("Wrong Room Id Provided : %v ",room_id)	
							
					server_response := ServerResponse{Type:"error",UserName:sender_name,Message:message,RoomId: room_id}
						
						conn.WriteJSON(server_response)
					  
						log.Println(message)

						continue
			}else {
				mu.Lock()
				users := room_map[joined_room_id]
				mu.Unlock()
				
				idx_to_delete := -1  
				for i, user := range users {
					if user.User_conn == conn {
						idx_to_delete = i
						break;
					}
				}
			
				if idx_to_delete != -1 {
					users = append(users[:idx_to_delete], users[idx_to_delete+1:]...)
				} 

				mu.Lock()
				room_map[joined_room_id] = users
				mu.Unlock()
       
				message := fmt.Sprintf("User %v left room %v",sender_name,joined_room_id)

				server_response := ServerResponse{Type:"leave",UserName:sender_name,Message:message,RoomId: room_id}

				for _,user := range users {
					user.User_conn.WriteJSON(server_response)
				}

				mu.Lock()
				delete(conn_map,conn)
				mu.Unlock()

				log.Printf("User %v left room %v",sender_name,joined_room_id)


			}

		}
		default : 
				if err := conn.WriteMessage(messageType,[]byte("Invalid Input ")) ; err != nil {
					log.Println(err)
				}
			

		}

	}
}


func ServerTester(){
 	
	for i:=0; i<1000; i++ {
		

		go func(i int ){
		fmt.Println("Writting to room : ", i)
				mu.Lock()
				room_map[1] = append(room_map[1],User{});			
				mu.Unlock()
		}(i)
	}

	fmt.Println("All Users Created Successfully ")

}


func main() {

	PORT := ":8080"

	log.Printf("Starting Server on  %v\n", PORT)

	go ServerTester();

	http.HandleFunc("/ws", MainHanlder)
	err := http.ListenAndServe(PORT, nil)

	if err != nil {
		log.Fatal(err)
	}

}
