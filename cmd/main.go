package main

import (
	"fmt"
)

func main() {
	options := []string{" Create Room ", " Join Room ", " Exit "}
	
	
	for i ,option := range options {
		fmt.Printf(" [ %d ] %s \n", i+1, option)
	}
	
	var input int
	fmt.Printf(" \n Select an option : ") 
	fmt.Scanln(&input)

	fmt.Println(options[input-1])

}