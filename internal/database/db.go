package database 

import (
	"fmt"
	"database/sql"
	_ "github.com/lib/pq"
)

const (
	host  = "localhost"
	port 	=  5432
	user 	= "postgres"
	password = "mypass"
	dbname = "termichatt"
)

func ConnectDatabse() (*sql.DB, error)  {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s" + 
	" password=%s dbname=%s sslmode=disable",host,port,user,password,dbname)
	


	db,err := sql.Open("postgres",psqlInfo)

	if err != nil {
		fmt.Println("Error Occured ",err)
		return nil,err
	}




	
	

	return db,nil

	



	
}




