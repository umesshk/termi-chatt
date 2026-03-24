package database

import (
	"log"
	"database/sql"
	"fmt"
)


func InsertUser(db *sql.DB, username string){
		
	log.Println("Inserting into users...")

	_, err := db.Exec("INSERT INTO USERS (username) VALUES ($1) ON CONFLICT (username) DO NOTHING ",username)

	if err != nil {
		log.Println("Error Inserting user ",err)
		return 
	}

 log.Printf("%v inserted in database... ", username)

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

func GetUserFromDB(db *sql.DB, username string ){

	var id int 

	query := fmt.Sprintf("SELECT id  FROM USERS WHERE USERSNAME=$1")

	err := db.QueryRow(query,username).Scan(&id)

	if err == nil {
		fmt.Printf("user: %v  not in database inserting...\n",username)
		InsertUser(db,username)
		return
	}

	fmt.Printf("user : %v  present in database ...\n",username)
	return 


}


