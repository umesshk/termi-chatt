package ws

import (
	"github.com/gorilla/websocket"
	"log"
	"fmt"
	"database/sql"
	"github.com/umesshk/termi-chatt/internal/database"
	userType  "github.com/umesshk/termi-chatt/internal/user"
	"context"
	"encoding/json"
	"time"
)



func HandleCreate(ClientMessage userType.UserMessage , conn *websocket.Conn,db *sql.DB, hub *Hub ){
			
				
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
				
				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				_ = hub.MarkRoomExists(ctx, roomId)
				cancel()

				user := userType.User{UserId: userId, Username: user_name, User_conn: conn}
				hub.AddConn(roomId, user)
				hub.EnsureRoomSub(roomId)


				message := fmt.Sprintf(" Created Room with Room Id :  %v",roomId)

				server_response := userType.ServerResponse{Type:"room_created", UserName:user_name, Message: message, RoomId: roomId}
			
				_ = conn.WriteJSON(server_response)

				log.Printf("%v  Created room %v\n" ,  user_name, roomId)
	
}


func HandleJoin(ClientMessage userType.UserMessage, conn *websocket.Conn,db *sql.DB, hub *Hub ){
	
				room_id   := ClientMessage.RoomId
				user_name := ClientMessage.Username
   			 
				userId,err :=  database.GetORInsertUser(db,user_name)
				
				user := userType.User{UserId: userId, Username: user_name, User_conn: conn}
				

				 if err != nil {
					 log.Println("Error Occured Inerting/getting user ",err)
					 return 
				 }

				log.Println("Client Room Id ", room_id)
				
				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				ok, existsErr := hub.RoomExists(ctx, room_id)
				cancel()
				if existsErr != nil {
					log.Println("room existence check:", existsErr)
					ok = false
				}

				
				if !ok {

					message := fmt.Sprintf("Room Doesn't Exist with room Id : %v", room_id)	
					server_response := userType.ServerResponse{Type:"error",UserName:user_name,Message:message,RoomId: room_id}
						
					_ = conn.WriteJSON(server_response)
					fmt.Println(message)
					return 

					}else {
							conn_room_id , ok := hub.JoinedRoomID(conn)

							if ok && conn_room_id == room_id {

									message := fmt.Sprintf("User %v Already  in  room Id : %v",user_name, room_id)	
						
									server_response := userType.ServerResponse{Type:"error",UserName:user_name,Message:message,RoomId: room_id}
					
									_ = conn.WriteJSON(server_response)
									fmt.Println(message)

						}else {

						hub.AddConn(room_id, user)
						hub.EnsureRoomSub(room_id)
			
						room_conn := hub.RoomUsers(room_id)
					
						database.UserJoinRoom(db,userId,room_id) 	
						
						var room_messages []userType.MessagesStruct

						message := fmt.Sprintf("%v  Joined the room ", user_name)
						
						room_messages, err = getRoomMessagesCached(db, hub, room_id)
						
						if err != nil {
							fmt.Println("Error Retreiving Messages : ", err)
							return 
						} 

						fmt.Println("Room Messages ", room_messages)
					
						server_response := userType.ServerResponse{Type:"room_joined",UserName:user_name,Message:message,RoomId: room_id}

							for _,receiver:= range room_conn {
						 		 
								receiver_conn := receiver.User_conn
					 		
								_ = receiver_conn.WriteJSON(server_response)


				 			log.Println(message)
						}
					  log.Println("Writting room messages to user ")	
						for _, msg := range room_messages {
							
								room_username := msg.Username
								message_content := msg.Content 
								server_response := userType.ServerResponse{Type:"chat_message",UserName:room_username,Message:message_content,RoomId: room_id}
							
								_ = conn.WriteJSON(server_response)
						}

					  log.Println(" room messages written to  user ")	
				}
			}
}

func HandleMessage( ClientMessage userType.UserMessage, conn *websocket.Conn,db *sql.DB, hub *Hub ){
 
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
						
						_ = conn.WriteJSON(server_response)
					  
						log.Println(message)

					return	
			
				}

				

				joined_room_id , ok := hub.JoinedRoomID(conn)

				if !ok {
					
					message := fmt.Sprintf("Please Join the Room First : %v ",room_id)	
							
						server_response := userType.ServerResponse{Type:"error",UserName:sender_name,Message:message,RoomId: room_id}
						
						_ = conn.WriteJSON(server_response)
					  
						log.Println(message)

					return	

			}else {

				if joined_room_id != room_id {
					message := fmt.Sprintf("Wrong Room Id Provided : %v ",room_id)	
							
					server_response := userType.ServerResponse{Type:"error",UserName:sender_name,Message:message,RoomId: room_id}
						
						_ = conn.WriteJSON(server_response)
					  
						log.Println(message)

					return	
			}

			 	 sender_message := ClientMessage.Message

				 log.Printf("Sender Message %v\n", sender_message)
					
				 database.InsertMessage(db,sender_id,joined_room_id,sender_message)
				 cacheAppendRoomMessage(hub, joined_room_id, sender_name, sender_message)
									
				 message_to_send := fmt.Sprintf("%v",sender_message)
				 server_response := userType.ServerResponse{Type:"chat_message", UserName:sender_name,Message:message_to_send,RoomId:room_id }
				 hub.Publish(joined_room_id, server_response)

			}

}

func HandleLeave( ClientMessage userType.UserMessage, conn *websocket.Conn ,db *sql.DB, hub *Hub){
				
				room_id := ClientMessage.RoomId
			  sender_name := ClientMessage.Username	
			
				
				if room_id == 0 {
						message := fmt.Sprintf("No Room Id provided... ")	
							
						server_response := userType.ServerResponse{Type:"error",UserName:sender_name,Message:message,RoomId: room_id}
						
						_ = conn.WriteJSON(server_response)
					  
						log.Println(message)

					return	
			
				}

				
				joined_room_id , ok := hub.JoinedRoomID(conn)

				if !ok {
					
					message := fmt.Sprintf("Please Join the Room First : %v ",room_id)	
							
						server_response := userType.ServerResponse{Type:"error",UserName:sender_name,Message:message,RoomId: room_id}
						
						_ = conn.WriteJSON(server_response)
					  
						log.Println(message)

					return	

			}else {

				if joined_room_id != room_id {
					message := fmt.Sprintf("Wrong Room Id Provided : %v ",room_id)	
							
					server_response := userType.ServerResponse{Type:"error",UserName:sender_name,Message:message,RoomId: room_id}
						
						_ = conn.WriteJSON(server_response)
					  
						log.Println(message)

					return	
			}else {
				hub.RemoveConn(conn)
				users := hub.RoomUsers(joined_room_id)
       
				message := fmt.Sprintf("User %v left room %v",sender_name,joined_room_id)

				server_response := userType.ServerResponse{Type:"leave",UserName:sender_name,Message:message,RoomId: room_id}

				for _,user := range users {
					_ = user.User_conn.WriteJSON(server_response)
				}

				log.Printf("User %v left room %v",sender_name,joined_room_id)


			}

		}
}

type cachedMsg struct {
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}

func getRoomMessagesCached(db *sql.DB, hub *Hub, roomID int) ([]userType.MessagesStruct, error) {
	if hub.redis == nil {
		return database.GetRoomMessages(db, roomID)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	key := fmt.Sprintf("room:%d:messages", roomID)
	vals, err := hub.redis.Rdb.LRange(ctx, key, 0, 49).Result()
	if err == nil && len(vals) > 0 {
		out := make([]userType.MessagesStruct, 0, len(vals))
		for i := len(vals) - 1; i >= 0; i-- { 	
			var cm cachedMsg
			if json.Unmarshal([]byte(vals[i]), &cm) == nil {
				out = append(out, userType.MessagesStruct{
					Username:  cm.Username,
					Content:   cm.Content,
					CreatedAt: cm.CreatedAt,
				})
			}
		}
		return out, nil
	}

	msgs, err := database.GetRoomMessages(db, roomID)
	if err != nil {
		return nil, err
	}

	pipe := hub.redis.Rdb.Pipeline()
	for i := len(msgs) - 1; i >= 0; i-- { 		
		b, _ := json.Marshal(cachedMsg{Username: msgs[i].Username, Content: msgs[i].Content, CreatedAt: msgs[i].CreatedAt})
		pipe.LPush(ctx, key, string(b))
	}
	pipe.LTrim(ctx, key, 0, 49)
	pipe.Expire(ctx, key, 10*time.Minute)
	_, _ = pipe.Exec(ctx)

	return msgs, nil
}

func cacheAppendRoomMessage(hub *Hub, roomID int, username, content string) {
	if hub.redis == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	key := fmt.Sprintf("room:%d:messages", roomID)
	b, err := json.Marshal(cachedMsg{Username: username, Content: content, CreatedAt: time.Now()})
	if err != nil {
		return
	}
	pipe := hub.redis.Rdb.Pipeline()
	pipe.LPush(ctx, key, string(b))
	pipe.LTrim(ctx, key, 0, 49)
	pipe.Expire(ctx, key, 10*time.Minute)
	_, _ = pipe.Exec(ctx)
}
