package database

import (
	"log"
	"database/sql"
	"fmt"
	"github.com/umesshk/termi-chatt/internal/user"
)


func GetORInsertUser(db *sql.DB, username string)(int,error){
		
	log.Println("Inserting into users...")
	
	var id int 

	err := db.QueryRow("INSERT INTO USERS (username) VALUES ($1) ON CONFLICT (username) DO UPDATE SET username=EXCLUDED.username RETURNING id",username).Scan(&id)

	if err != nil {
		return 0,err
	}
 
log.Printf("User %v has id %v", username, id)
 
 return id, nil  


}

func CreateRoom(db *sql.DB  )(int,error){
	
	log.Println("Creating room database...")
	
	query := fmt.Sprintf("INSERT INTO ROOMS DEFAULT VALUES RETURNING id")
	
	var room_id int 

	 err := db.QueryRow(query).Scan(&room_id)
	
	 if err != nil {
		 return 0, err
	 }
	 return room_id,nil

}

func UserJoinRoom(db *sql.DB,  user_id,room_id  int ){
   
	log.Println("user joins room ")	
	
	query := fmt.Sprintf("INSERT INTO room_users (user_id, room_id) VALUES ($1,$2) ")


	_,err := db.Exec(query,user_id,room_id)

	if err != nil {
		log.Println("Error user joining room ",err)
		return 
	} 

	log.Printf("User %v inserted with room %v in database", user_id,room_id )

}


func InsertMessage(db *sql.DB , user_id , room_id int, message string){

log.Println("Inserting Messages into Database ")
 query := fmt.Sprintf("INSERT INTO MESSAGES (user_id, room_id,content) VALUES ($1,$2,$3)")	

 _,err := db.Exec(query,user_id,room_id,message)

 if err != nil {
	 log.Println("Error Inserting Messages ",err)
	 return 
 }

 log.Printf("User %v inserted in room %v message %v", user_id,room_id,message)

}

func GetRoomMessages(db *sql.DB, room_id int ) (user.Messages[], error){

	log.Println("Getting Room %v Messages from Database")

	rows, err := db.Query(
		`SELECT u.username , m.content, m.created_at 
		 FROM messages m 
		 JOIN users u on m.user_id = u.id 
		 WHERE m.room_id = $1 
		 ORDER BY m.created_at ASC 
		 LIMIT 50
		`
	,room_id)

  if err != nil {
		return _, err
	}

	var messages []user.Messages

	defer rows.Close()

	for rows.Next(){
		var msg user.Messages


		err := rows.Scan(&msg.Username , &msg.Content, &msg.CreatedAt)

		if err != nil {
			return _,err
		}


		messages = append(messages,msg)
		


	}
   return messages, nil

}



