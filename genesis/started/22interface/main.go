package main

import "fmt"

func main() {

	dog := Dog{
		Name: "Beko",
		Age:  3,
	}

	cow := Cow{
		Name: "kyolaba",
		Age:  5,
	}

	PrintInfo(&dog)
	PrintInfo(&cow)

}

// define nethods and their return type
type Animal interface {
	Says() string
	NumLegs() int
}

func PrintInfo(a Animal) {
	fmt.Println("This amimal says", a.Says(), "and has num of legs ", a.NumLegs())
}

type Dog struct {
	Name string
	Age  int
}

type Cow struct {
	Name string
	Age  int
}

func (d *Dog) Says() string {
	return "bobo"
}
func (d *Dog) NumLegs() int {
	return 4
}

func (d *Cow) Says() string {
	return "mooo"
}
func (d *Cow) NumLegs() int {
	return 4
}
