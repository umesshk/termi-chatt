package server 

import (
		"log"
		"net/http"
)
func StartServer(){ 
	
	PORT := ":8080"
	
	server := &http.Server{
		Addr : PORT,
		Handler: nil, 
	}

	log.Printf("Starting Server on  %v\n", PORT)
	
	err := server.ListenAndServe()

	if err!=nil {
		log.Fatal(err)
	}

}
