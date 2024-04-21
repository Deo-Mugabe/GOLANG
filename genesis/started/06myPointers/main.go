package main

import "fmt"

func main() {
	fmt.Println("Welcome to pointers")

	myvalue := 23

	var ptr = &myvalue

	fmt.Println("my pointer address is ", ptr)
	fmt.Println("my pointer value is ", *ptr)

	*ptr = *ptr + 2

	fmt.Println("my pointer value is ", myvalue)
}
