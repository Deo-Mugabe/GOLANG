package main

import "fmt"

func main() {

	John := User{"John", 23, true}
	fmt.Println(John)
	fmt.Println(John.Age)
	John.GetStatus()
	John.SetAge()
}

type User struct {
	Name   string
	Age    int
	Status bool
}

func (u User) GetStatus() {
	fmt.Println(u.Status)
}

func (u User) SetAge() {
	u.Age = 34
	fmt.Println(u.Age)
}
