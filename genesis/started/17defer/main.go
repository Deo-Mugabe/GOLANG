package main

import "fmt"

func main() {
	defer fmt.Println("World")
	defer fmt.Println("One")
	defer fmt.Println("Two")
	fmt.Println("Hello")
	myDefer()

	fmt.Println("Start")
	panic("This is a panic the nxt statement wont be executed") // unrecoverable condition but u can recover it
	fmt.Println("End")
}

func myDefer() {
	for i := 0; i < 5; i++ {
		defer fmt.Print(i)
	}
}
