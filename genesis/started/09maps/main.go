package main

import "fmt"

func main() {

	fmt.Println("Welcome to maps")

	myMap := make(map[string]string)

	myMap["py"] = "Python"
	myMap["JS"] = "JavaScript"
	myMap["GO"] = "GoLang"

	fmt.Println("my languages ", myMap)
	fmt.Println("my languages ", myMap["GO"])

	// delete
	delete(myMap, "JS")
	fmt.Println("my languages ", myMap)

	for key, value := range myMap {
		fmt.Printf("The Key %v, value is %v\n", key, value)
	}

	for _, value := range myMap {
		fmt.Printf("The value is %v\n", value)
	}

}
