package server

import (
	"log"
	"net"
)


func handleConnection(conn net.Conn){
		log.Printf("handle func")
}

func CreateServer(){
	ln, err := net.Listen("tcp", ":8080")
	log.Printf("Listening of port 8080")
	if err!=nil {
		log.Fatal("An Error Occured starting server")
	}

	for {
		conn, err := ln.Accept()
		if err!=nil{
				log.Fatal("Couldn't Connect")
		}

		go handleConnection(conn)
	}
	
}