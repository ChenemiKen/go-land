package main

import (
	"fmt"
	"log"
)

func main() {
	fmt.Println("Hello, world!")

	var whatToSay string
	var i int

	whatToSay = "Goodbye cruel world"
	fmt.Println(whatToSay)

	i = 8
	fmt.Println("i is set to", i)

	changeWithPointer(&whatToSay)
	log.Println("after function call", whatToSay)
}

func changeWithPointer(s *string) {
	newValue := "Red"
	*s = newValue
}
