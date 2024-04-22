package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	fmt.Println("Welcome to our game")

	rand.Seed(time.Now().UnixNano())
	diceNumber := rand.Intn(6) + 1

	switch diceNumber {

	case 1:
		fmt.Print("Dice value is 1 and you can open")

	case 2:
		fmt.Print("Dice value is 2 spot")

	case 3:
		fmt.Print("Dice value is 3 spot")

	case 4:
		fmt.Print("Dice value is 4 spot")

	case 5:
		fmt.Print("Dice value is 5 spot")

	case 6:
		fmt.Print("Dice value is 6 roll dice again")

	default:
		fmt.Println("what app")
	}
}
