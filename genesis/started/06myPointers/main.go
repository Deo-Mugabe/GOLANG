package main

import (
	"fmt"
	"log"
)

func main() {
	fmt.Println("Welcome to pointers")

	myvalue := 23

	var ptr = &myvalue

	fmt.Println("my pointer address is ", ptr)
	fmt.Println("my pointer value is ", *ptr)

	*ptr = *ptr + 2

	fmt.Println("my pointer value is ", myvalue)

	var myString string
	myString = "Green"
	log.Println("Before func call my string is ", myString)
	changeUsingPointer(&myString)
	log.Println("after func call string is set to ", myString)
}

func changeUsingPointer(s *string) {
	log.Println("s is set to ", s)
	newValue := "Red"
	*s = newValue
}
