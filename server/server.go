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

	log.Printf("Server Running on  %v", PORT)
	
	err := server.ListenAndServe()

	if err!=nil {
		log.Fatal(err)
	}

}
