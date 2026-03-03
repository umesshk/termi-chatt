package main

import (
	"fmt"
	"os"
	"github.com/umeshhk/termi-chatt/server"
)

func CreateRoom(){
	fmt.Println("Room Created")
}

func JoinRoom(){
	fmt.Println("Room Joined")
}

func main() {
	options := []string{" Create Room ", " Join Room ", " Exit "}
	
	
	for i ,option := range options {
		fmt.Printf(" [ %d ] %s \n", i+1, option)
	}
	
	var user_choice int
	fmt.Printf(" \n Select an option : ") 
	fmt.Scanln(&user_choice)
 
	switch user_choice { 
	case 1 : 
	 	fmt.Println("Creating Room ")
		CreateRoom()
 	
	case 2: 
		fmt.Println("Joining Room ")
 		JoinRoom()
	 
	case 3:
		fmt.Println("Exiting Program...")
		os.Exit(0)

	default  : 
		fmt.Println("Invalid Choice ")
}
}
