package main

import (
	"encoding/json"
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
}
