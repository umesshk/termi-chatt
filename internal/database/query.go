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



