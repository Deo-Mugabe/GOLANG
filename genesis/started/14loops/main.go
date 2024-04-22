package main

import "fmt"

func main() {
	fmt.Println("Welcome to loops")

	days := []string{"Monday", "Tuesday", "Wednesday"}

	for i, day := range days {

		fmt.Printf("Index %v of day %v\n", i, day)
	}

	myval := 1

	for myval < 10 {

		if myval == 5 {
			goto lco
		}
		fmt.Println(myval)
		myval++
	}

lco:
	fmt.Println("Hello ")

}
