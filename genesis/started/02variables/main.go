package main

import (
	"fmt"
	"strconv"
)

// constant
const LoginToken string = "dffeagad" /// the name having the capital letter means its public

var i = 32

func main() {

	i := 10 // by this am shadowing the globle scope i
	fmt.Println(i)

	var j float32 = float32(i) // type casting
	fmt.Println(j)

	k := 23 // by default Go will make it float64
	fmt.Printf("%v, %T", k, k)

	var toString string = strconv.Itoa(i) // lconverting to string
	fmt.Println(toString)

	var username string = "Jesus"
	fmt.Println(username)
	fmt.Printf("Variable is of type: %T \n", username) // %T is a place holder

	var anotherVal2 string
	fmt.Println(anotherVal2)
	fmt.Printf("variable is of type: %T \n", anotherVal2)

	var website = "Jesus is Lord" // by default lexer will make it string base on ur data
	fmt.Println(website)
	//no var sytle
	numberOfUser, num2 := 43434, 343 /// := can not be used ouside the function
	fmt.Println(numberOfUser, num2)
}
