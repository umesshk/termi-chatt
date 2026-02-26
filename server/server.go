package server 

import (
		"log"
		"net/http"
)

func HomeHandler(w http.ResponseWriter,r *http.Request){
	bytes, err := w.Write([]byte("Hello from Server "))

	if err != nil {
		log.Fatal(err)
		return
	}
	log.Fatal("%v Bytes of Data written Succesfully ",bytes)
} 


func StartServer(){ 
	
	PORT := ":8080"
	

	log.Printf("Starting Server on  %v\n", PORT)
	
	err := http.ListenAndServe(PORT,nil)

	if err!=nil {
		log.Fatal(err)
	}

	http.HandleFunc("/hello",HomeHandler)

}
