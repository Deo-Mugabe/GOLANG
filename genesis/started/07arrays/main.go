package main

import "fmt"

func main() {
	fmt.Println("Welcome to arrays in GO")

	var fruitList [4]string
	fruitList[0] = "apple"
	fruitList[1] = "mango"
	fruitList[2] = "pine"
	fruitList[3] = "cone"

	fmt.Println("List of fruits: ", fruitList)

	var vegList = [3]string{"dodo", "spinach", "pine"}

	fmt.Println("List of fruits: ", vegList)

}
