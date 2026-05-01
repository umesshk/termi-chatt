package database 

import (
	"fmt"
	"database/sql"
	_ "github.com/lib/pq"
)

func ConnectDatabse(dsn string) (*sql.DB, error) {
	if dsn == "" {
		return nil, fmt.Errorf("empty postgres dsn")
	}

	db,err := sql.Open("postgres", dsn)

	if err != nil {
		fmt.Println("Error Occured ",err)
		return nil,err
	}

	return db,nil


	
}




