package ws

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/umesshk/termi-chatt/internal/database"
	userType "github.com/umesshk/termi-chatt/internal/user"
)

func HandleCreate(clientMessage userType.UserMessage, client *Client, db *sql.DB, hub *Hub) {
	userName := clientMessage.Username

	userID, err := database.GetORInsertUser(db, userName)
	if err != nil {
		log.Println("Error inserting user  ", err)
		return
	}

	roomID, err := database.CreateRoom(db)
	if err != nil {
		log.Println("Error Creating room ", err)
		return
	}

	log.Println("Room created with id ", roomID)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	_ = hub.MarkRoomExists(ctx, roomID)
	cancel()

	hub.AddClient(roomID, userID, userName, client)
	hub.EnsureRoomSub(roomID)

	message := fmt.Sprintf(" Created Room with Room Id :  %v", roomID)
	serverResponse := userType.ServerResponse{Type: "room_created", UserName: userName, Message: message, RoomId: roomID}

	client.Enqueue(serverResponse)

	log.Printf("%v  Created room %v\n", userName, roomID)
}

func HandleJoin(clientMessage userType.UserMessage, client *Client, db *sql.DB, hub *Hub) {
	roomID := clientMessage.RoomId
	userName := clientMessage.Username

	userID, err := database.GetORInsertUser(db, userName)
	if err != nil {
		log.Println("Error Occured Inerting/getting user ", err)
		return
	}

	log.Println("Client Room Id ", roomID)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	ok, existsErr := hub.RoomExists(ctx, roomID)
	cancel()
	if existsErr != nil {
		log.Println("room existence check:", existsErr)
		ok = false
	}

	if !ok {
		message := fmt.Sprintf("Room Doesn't Exist with room Id : %v", roomID)
		serverResponse := userType.ServerResponse{Type: "error", UserName: userName, Message: message, RoomId: roomID}

		client.Enqueue(serverResponse)
		fmt.Println(message)
		return
	}

	connRoomID, ok := hub.JoinedRoomID(client)
	if ok && connRoomID == roomID {
		message := fmt.Sprintf("User %v Already  in  room Id : %v", userName, roomID)
		serverResponse := userType.ServerResponse{Type: "error", UserName: userName, Message: message, RoomId: roomID}

		client.Enqueue(serverResponse)
		fmt.Println(message)
		return
	}

	hub.AddClient(roomID, userID, userName, client)
	hub.EnsureRoomSub(roomID)

	database.UserJoinRoom(db, userID, roomID)

	var roomMessages []userType.MessagesStruct

	message := fmt.Sprintf("%v  Joined the room ", userName)

	roomMessages, err = getRoomMessagesCached(db, hub, roomID)
	if err != nil {
		fmt.Println("Error Retreiving Messages : ", err)
		return
	}

	fmt.Println("Room Messages ", roomMessages)

	serverResponse := userType.ServerResponse{Type: "room_joined", UserName: userName, Message: message, RoomId: roomID}
	hub.BroadcastToRoom(roomID, serverResponse)

	log.Println(message)
	log.Println("Writting room messages to user ")
	for _, msg := range roomMessages {
		roomUsername := msg.Username
		messageContent := msg.Content
		historyResponse := userType.ServerResponse{Type: "chat_message", UserName: roomUsername, Message: messageContent, RoomId: roomID}

		client.Enqueue(historyResponse)
	}

	log.Println(" room messages written to  user ")
}

func HandleMessage(clientMessage userType.UserMessage, client *Client, db *sql.DB, hub *Hub) {
	roomID := clientMessage.RoomId
	senderName := clientMessage.Username

	senderID, err := database.GetORInsertUser(db, senderName)
	if err != nil {
		log.Println("Error Occured getting user id ", err)
		return
	}

	if roomID == 0 {
		message := fmt.Sprintf("No Room Id provided... ")
		serverResponse := userType.ServerResponse{Type: "error", UserName: senderName, Message: message, RoomId: roomID}

		client.Enqueue(serverResponse)
		log.Println(message)
		return
	}

	joinedRoomID, ok := hub.JoinedRoomID(client)
	if !ok {
		message := fmt.Sprintf("Please Join the Room First : %v ", roomID)
		serverResponse := userType.ServerResponse{Type: "error", UserName: senderName, Message: message, RoomId: roomID}

		client.Enqueue(serverResponse)
		log.Println(message)
		return
	}

	if joinedRoomID != roomID {
		message := fmt.Sprintf("Wrong Room Id Provided : %v ", roomID)
		serverResponse := userType.ServerResponse{Type: "error", UserName: senderName, Message: message, RoomId: roomID}

		client.Enqueue(serverResponse)
		log.Println(message)
		return
	}

	senderMessage := clientMessage.Message

	log.Printf("Sender Message %v\n", senderMessage)

	database.InsertMessage(db, senderID, joinedRoomID, senderMessage)
	cacheAppendRoomMessage(hub, joinedRoomID, senderName, senderMessage)

	messageToSend := fmt.Sprintf("%v", senderMessage)
	serverResponse := userType.ServerResponse{Type: "chat_message", UserName: senderName, Message: messageToSend, RoomId: roomID}
	hub.Publish(joinedRoomID, serverResponse)
}

func HandleLeave(clientMessage userType.UserMessage, client *Client, db *sql.DB, hub *Hub) {
	roomID := clientMessage.RoomId
	senderName := clientMessage.Username

	if roomID == 0 {
		message := fmt.Sprintf("No Room Id provided... ")
		serverResponse := userType.ServerResponse{Type: "error", UserName: senderName, Message: message, RoomId: roomID}

		client.Enqueue(serverResponse)
		log.Println(message)
		return
	}

	joinedRoomID, ok := hub.JoinedRoomID(client)
	if !ok {
		message := fmt.Sprintf("Please Join the Room First : %v ", roomID)
		serverResponse := userType.ServerResponse{Type: "error", UserName: senderName, Message: message, RoomId: roomID}

		client.Enqueue(serverResponse)
		log.Println(message)
		return
	}

	if joinedRoomID != roomID {
		message := fmt.Sprintf("Wrong Room Id Provided : %v ", roomID)
		serverResponse := userType.ServerResponse{Type: "error", UserName: senderName, Message: message, RoomId: roomID}

		client.Enqueue(serverResponse)
		log.Println(message)
		return
	}

	hub.RemoveClient(client)

	message := fmt.Sprintf("User %v left room %v", senderName, joinedRoomID)
	serverResponse := userType.ServerResponse{Type: "leave", UserName: senderName, Message: message, RoomId: roomID}
	hub.BroadcastToRoom(joinedRoomID, serverResponse)

	log.Printf("User %v left room %v", senderName, joinedRoomID)
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
