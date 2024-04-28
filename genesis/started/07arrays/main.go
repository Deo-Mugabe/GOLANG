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

	// you have to declare the size of the array but if you dont want to
	//define the size go but data after initializing

	var myArray = [...]int{1, 2, 3, 4, 5, 6}
	fmt.Println(myArray[3])

	b := &myArray // by this we are not creating a new array so change of b means original is changed

	b[0] = 10
	fmt.Println(myArray)

	c := myArray[4:]
	fmt.Println(c)

	d := myArray[:len(myArray)-1] // remove the last element in an array
	fmt.Println(d)

	// multi dimensional array

	var multiD [3][3]int = [3][3]int{
		{1, 2, 5},
		{0, 2, 0},
		{1, 3, 3},
	}
	fmt.Println(multiD)
}
