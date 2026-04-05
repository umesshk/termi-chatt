
package user

import (
	"github.com/gorilla/websocket"
)

type UserMessage struct {
 Msgtype 		string  `json:"type"` 
 Message		string  `json:"message,omitempty"` 
 RoomId   	int			`json:"roomId,omitempty"`
 Username 	string 	`json:"username"`
}

type User struct {
	UserId   		int 						`json:"userId"`
	Username 		string 					`json:"username"`
  User_conn  *websocket.Conn 	`json:"conn,omitempty"`
}

type ServerResponse struct {
	Type   			string 	`json:"type"`
	UserName 		string 	`json:"username"`
	Message     string  `json:"message"`
	RoomId			int 		`json:"roomId"`
}

type MessagesStruct struct {
	Username	 string   
	Content 	 string  
	CreatedAt  string   
}
