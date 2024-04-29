package main

import (
	"fmt"
	"log"
)

// collection of objects of different datatypes

type User struct {
	Name   string
	Email  string
	Status bool
	Age    int
}

type Student struct {
	Name   string
	Email  string
	Status bool
	Age    int
}

func main() {

	Student1 := Student{Name: "Deo", Status: true}
	Student2 := Student1
	fmt.Println(Student1)
	fmt.Println(Student2)
	Student2.Name = "Eric"
	Student3 := &Student1
	Student3.Name = "Degra"
	fmt.Println(Student1)
	fmt.Println(Student2)
	fmt.Println(Student3)

	fmt.Println("welcome to Struts which are the classes in java")

	myUser := User{"Deo", "deo@go.dev", true, 33}

	fmt.Printf("User profile is %+v\n", myUser)

	fmt.Printf("The name of is %v and email is %v", myUser.Name, myUser.Email)
	log.Println("Print name using method ", myUser.printName())

}

func (m *User) printName() string {
	return m.Name
}
