package server

import (
	"fmt"
	"log"
	"net/http"
)

func Handlehttp(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello from the server ")
}

func StartServer() {
	http.HandleFunc("/", Handlehttp)
	log.Println("Server is running at port :8080")
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Println("An Error Occured starting server")
	}
}
