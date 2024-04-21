package main

import "fmt"

//constant
const LoginToken string = "dffeagad" /// the name having the capital letter means its public

func main() {
	var username string = "Jesus"
	fmt.Println(username)
	fmt.Printf("Variable is of type: %T \n", username) // %T is a place holder

	var isLoggedin bool = true
	fmt.Println(isLoggedin)
	fmt.Printf("varriable is of type: %T \n", isLoggedin)

	var smallval uint8 = 255
	fmt.Println(smallval)
	fmt.Printf("varriable is of type: %T \n", smallval)

	var anotherVal int
	fmt.Println(anotherVal)
	fmt.Printf("variable is of type: %T \n", anotherVal)

	var anotherVal2 string
	fmt.Println(anotherVal2)
	fmt.Printf("variable is of type: %T \n", anotherVal2)

	var website = "Jesus is Lord" // by default lexer will make it string base on ur data
	fmt.Println(website)
	//no var sytle
	numberOfUser, num2 := 43434, 343 /// := can not be used ouside the function
	fmt.Println(numberOfUser, num2)
}
