package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type Person struct {
	First_Name string `json:"first_name"`
	Last_Name  string `json:"last_name"`
	IsMarried  bool   `json:"isMarried"`
}

func main() {
	myJson := `
	[
		{
			"first_name": "Deo",
			"last_name": "Mugabe",
			"isMarried": true
		},
		{
			"first_name": "Eric",
			"last_name": "Kikomeko",
			"isMarried": true
		}
	]`

	var unmarshalled []Person

	err := json.Unmarshal([]byte(myJson), &unmarshalled)
	if err != nil {
		log.Println("Error while Umarshalling", err)
	}
	log.Printf("Unmarshalled: %v", unmarshalled)

	// Write json from a struct

	var mySlice []Person

	var m1 Person
	m1.First_Name = "John"
	m1.Last_Name = "Doe"
	m1.IsMarried = false

	var m2 Person
	m2.First_Name = "John"
	m2.Last_Name = "Doe"
	m2.IsMarried = false

	mySlice = append(mySlice, m2)

	newJson, err := json.MarshalIndent(mySlice, "", "	")
	if err != nil {
		log.Println("error marshalling", err)
	}

	fmt.Println(string(newJson))

}
