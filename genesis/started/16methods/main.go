package main

import (
	"fmt"
	"math"
)

// Methods work on type define datatypes ( struts)
func main() {

	John := User{"John", 23, true}
	fmt.Println(John)
	fmt.Println(John.Age)
	John.GetStatus()
	John.SetAge()

	var r geometry = rectangle{height: 2.0, width: 4.0}
	fmt.Println("Area of a rectangle: ", r.area())
	fmt.Println("Perimeter of a rectangle: ", r.perimeter())

	var c geometry = cicle{radis: 2.0}
	fmt.Println("Area of a circle: ", c.area())
	fmt.Println("Perimeter of a circle: ", c.perimeter())
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

type geometry interface {
	area() float64
	perimeter() float64
}

type cicle struct {
	radis float64
}

type rectangle struct {
	height, width float64
}

func (r rectangle) area() float64 {
	return r.height * r.width
}

func (r rectangle) perimeter() float64 {
	return 2*r.height + 2*r.width
}

func (c cicle) area() float64 {
	return math.Pi * c.radis * c.radis
}

func (c cicle) perimeter() float64 {
	return 2 * math.Pi * c.radis
}
