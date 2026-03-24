package database

import (
	"log"
	"database/sql"
	"fmt"
)


func InsertUser(db *sql.DB, username string)(int,error){
		
	log.Println("Inserting into users...")
	
	var id int 

	 err := db.QueryRow("INSERT INTO USERS (username) VALUES ($1) ON CONFLICT (username) DO UPDATE SET username=EXCLUDED.username RETURNING id",username).Scan(&id)

	if err != nil {
		return 0,err
	}

 log.Printf("%v inserted in database... ", username)

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





