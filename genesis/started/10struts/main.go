package main

import (
	"fmt"
	"log"
)

type User struct {
	Name   string
	Email  string
	Status bool
	Age    int
}

func main() {
	fmt.Println("welcome to Struts which are the classes in java")

	myUser := User{"Deo", "deo@go.dev", true, 33}

	fmt.Printf("User profile is %+v\n", myUser)

	fmt.Printf("The name of is %v and email is %v", myUser.Name, myUser.Email)
	log.Println("Print name using method ", myUser.printName())

}

func (m *User) printName() string {
	return m.Name
}
