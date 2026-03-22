package database

import (
	"log"
	"database/sql"
	"fmt"
)


func InsertUser(db *sql.DB, username string){
		
	fmt.Println("Inserting into users...")

	query := fmt.Sprintf("INSERT INTO USERS (username) VALUES ($1)")
	
 _ ,err := db.Exec(query,username)

 if err !=nil {
	 log.Panic("Error Inserting Value ",err)
	 return 
 }

 log.Printf("%v inserted in database... ", username)

}

func CreateRoom(db *sql.DB, room_id int ){
	
	log.Println("Inserting room to database")
	
	query := fmt.Sprintf("INSERT INTO ROOMS (room_id) VALUES ($1)")

	_, err := db.Exec(query,room_id)

	if err != nil {
		log.Panic("Error Inserting room to db", err)
		return 
	}

 log.Printf("%v room  inserted in database... ", room_id)

}



