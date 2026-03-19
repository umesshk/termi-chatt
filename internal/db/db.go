package db 

import (
	"fmt"
	"database/sql"
	_ "github.com/lib/pq"
)

const (
	host  = "localhost"
	port 	=  5432
	user 	= "postgres"
	password = "myspass"
	dbname = "termiChatt"
)

func ConnectDatabse()  {

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s" + "password=%s dbname=%s sslmode=disable",host,port,user,password,dbname)

	db,err := sql.Open("postgres",psqlInfo)

	if err != nil {
		fmt.Println("Error Occured ",err)
		return 
	}

	defer db.Close()

	fmt.Println("Connected to Database...")
	
}




