package ws 

import (
	"github.com/gorilla/websocket"
	"log"
	"fmt"
	"sync"
	"database/sql"
	"github.com/umesshk/termi-chatt/internal/database"
	userType  "github.com/umesshk/termi-chatt/internal/user"
)

var mu sync.Mutex


//  roomid -> websocket connection slice
var room_map = make(map[int][]userType.User)
// websocket connection -> room id 
var conn_map = make(map[*websocket.Conn]int)



func HandleCreate(ClientMessage userType.UserMessage , conn *websocket.Conn,db *sql.DB ){
			
				
				user_name := ClientMessage.Username
				
				userId,err := database.GetORInsertUser(db,user_name)
				
				if err != nil {
					log.Println("Error inserting user  ",err)
					return 
				}

				roomId,err := database.CreateRoom(db) 
			
				if err != nil {
					log.Println("Error Creating room ",err)
					return 
				}

				log.Println("Room created with id ", roomId)
				

				

				user := userType.User{userId,user_name,conn}
				
				mu.Lock()
				room_map[roomId] = append(room_map[roomId], user)	

	    	conn_map[conn] = roomId	

				mu.Unlock()


				message := fmt.Sprintf(" Created Room with Room Id :  %v",roomId)

				server_response := userType.ServerResponse{Type:"room_created", UserName:user_name, Message: message, RoomId: roomId}
			
				conn.WriteJSON(server_response)

				log.Printf("%v  Created room %v\n" ,  user_name, roomId)
	
}


func HandleJoin(ClientMessage userType.UserMessage, conn *websocket.Conn,db *sql.DB ){
	
				room_id   := ClientMessage.RoomId
				user_name := ClientMessage.Username
   			 
				userId,err :=  database.GetORInsertUser(db,user_name)
				
				user := userType.User{userId, user_name, conn}
				

				 if err != nil {
					 log.Println("Error Occured Inerting/getting user ",err)
					 return 
				 }

				log.Println("Client Room Id ", room_id)
				
				mu.Lock()
				_, ok := room_map[room_id]	
				mu.Unlock()

				
				if !ok {

					message := fmt.Sprintf("Room Doesn't Exist with room Id : %v", room_id)	
					server_response := userType.ServerResponse{Type:"error",UserName:user_name,Message:message,RoomId: room_id}
						
					conn.WriteJSON(server_response)
					fmt.Println(message)
					return 

					}else {
							mu.Lock()
					   	conn_room_id , ok := conn_map[conn]
							mu.Unlock()

							if ok && conn_room_id == room_id {

									message := fmt.Sprintf("User %v Already  in  room Id : %v",user_name, room_id)	
						
									server_response := userType.ServerResponse{Type:"error",UserName:user_name,Message:message,RoomId: room_id}
					
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
					
						database.UserJoinRoom(db,userId,room_id) 	
						
						var room_messages []userType.MessagesStruct

						message := fmt.Sprintf("%v  Joined the room ", user_name)
						
						room_messages , err := database.GetRoomMessages(db,room_id) 
						
						if err != nil {
							fmt.Println("Error Retreiving Messages : ", err)
							return 
						} 

						fmt.Println("Room Messages ", room_messages)
					
						server_response := userType.ServerResponse{Type:"room_joined",UserName:user_name,Message:message,RoomId: room_id}

							for _,receiver:= range room_conn {
						 		 
								receiver_conn := receiver.User_conn
					 		
								receiver_conn.WriteJSON(server_response)


				 			log.Println(message)
						}
					  log.Println("Writting room messages to user ")	
						for _, msg := range room_messages {
							
								room_username := msg.Username
								message_content := msg.Content 
								server_response := userType.ServerResponse{Type:"chat_message",UserName:room_username,Message:message_content,RoomId: room_id}
							
								conn.WriteJSON(server_response)
						}

					  log.Println(" room messages written to  user ")	
				}
			}
}

func HandleMessage( ClientMessage userType.UserMessage, conn *websocket.Conn,db *sql.DB ){
 
				room_id := ClientMessage.RoomId
			  sender_name := ClientMessage.Username	
				
				sender_id,err := database.GetORInsertUser(db,sender_name)

				if err != nil {
				log.Println("Error Occured getting user id ",err)
					return
				}



				if room_id == 0 {
						message := fmt.Sprintf("No Room Id provided... ")	
							
						server_response := userType.ServerResponse{Type:"error",UserName:sender_name,Message:message,RoomId: room_id}
						
						conn.WriteJSON(server_response)
					  
						log.Println(message)

					return	
			
				}

				

				mu.Lock()
				joined_room_id , ok := conn_map[conn]
				mu.Unlock()

				if !ok {
					
					message := fmt.Sprintf("Please Join the Room First : %v ",room_id)	
							
						server_response := userType.ServerResponse{Type:"error",UserName:sender_name,Message:message,RoomId: room_id}
						
						conn.WriteJSON(server_response)
					  
						log.Println(message)

					return	

			}else {

				if joined_room_id != room_id {
					message := fmt.Sprintf("Wrong Room Id Provided : %v ",room_id)	
							
					server_response := userType.ServerResponse{Type:"error",UserName:sender_name,Message:message,RoomId: room_id}
						
						conn.WriteJSON(server_response)
					  
						log.Println(message)

					return	
			}

				mu.Lock()
 				 current_room := room_map[joined_room_id]
				 mu.Unlock()
			 	 sender_message := ClientMessage.Message

				 log.Printf("Sender Message %v\n", sender_message)
					
				 database.InsertMessage(db,sender_id,joined_room_id,sender_message)
									
					fmt.Println("Writing Message to all user ")
			
					for _,receiver:= range current_room {

					 receiver_name := receiver.Username
					 receiver_conn := receiver.User_conn
					 
				 	 message_to_send := fmt.Sprintf("%v",sender_message)
					
					 server_response := userType.ServerResponse{Type:"chat_message", UserName:sender_name,Message:message_to_send,RoomId:room_id }
						
					 receiver_conn.WriteJSON(server_response)
					 
				 log.Printf("Message Written to User :  %v \n",receiver_name)
				 


				 }

			}

}

func HandleLeave( ClientMessage userType.UserMessage, conn *websocket.Conn ,db *sql.DB){
				
				room_id := ClientMessage.RoomId
			  sender_name := ClientMessage.Username	
			
				
				if room_id == 0 {
						message := fmt.Sprintf("No Room Id provided... ")	
							
						server_response := userType.ServerResponse{Type:"error",UserName:sender_name,Message:message,RoomId: room_id}
						
						conn.WriteJSON(server_response)
					  
						log.Println(message)

					return	
			
				}

				
				mu.Lock()
				joined_room_id , ok := conn_map[conn]
				mu.Unlock()

				if !ok {
					
					message := fmt.Sprintf("Please Join the Room First : %v ",room_id)	
							
						server_response := userType.ServerResponse{Type:"error",UserName:sender_name,Message:message,RoomId: room_id}
						
						conn.WriteJSON(server_response)
					  
						log.Println(message)

					return	

			}else {

				if joined_room_id != room_id {
					message := fmt.Sprintf("Wrong Room Id Provided : %v ",room_id)	
							
					server_response := userType.ServerResponse{Type:"error",UserName:sender_name,Message:message,RoomId: room_id}
						
						conn.WriteJSON(server_response)
					  
						log.Println(message)

					return	
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

				server_response := userType.ServerResponse{Type:"leave",UserName:sender_name,Message:message,RoomId: room_id}

				for _,user := range users {
					user.User_conn.WriteJSON(server_response)
				}

				mu.Lock()
				delete(conn_map,conn)
				mu.Unlock()

				log.Printf("User %v left room %v",sender_name,joined_room_id)


			}

		}
}
