package main

import (
	"fmt"
	"sort"
)

func main() {
	fmt.Println("Welcome to slices")

	var sampelSlice = []string{"dfd", "dfdf", "erer"}
	fmt.Println(sampelSlice)

	sampelSlice = append(sampelSlice, "kkk", "kkk3")
	fmt.Println(sampelSlice)

	sampelSlice = append(sampelSlice[1:2])
	fmt.Println(sampelSlice)

	highscores := make([]int, 4)
	highscores[0] = 2
	highscores[1] = 3
	highscores[2] = 5
	highscores[3] = 8
	fmt.Println(highscores)

	highscores = append(highscores, 23, 2, 43, 45)
	fmt.Println(highscores)

	sort.Ints(highscores)
	fmt.Println(highscores)

	var courses = []string{"sa", "ea", "AI", "DBMS", "FPP", "MPP"}

	fmt.Println(courses)

	var index int = 2
	courses = append(courses[:index], courses[index+1:]...)
	fmt.Println(courses)

	// slices work like arrays difference is slices dont have size and also if you declare a new array
	// from another array it creates a copy of that array, but for the slice it points to the same
	// memory address. (uses pointers internally)

	var a = []int{1, 2, 3, 4, 5}
	var b = a

	fmt.Println(a)
	fmt.Println(b)
	b[0] = 9
	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(len(a)) // len and cap are different even though may produce same results
	fmt.Println(cap(b))
	c := append(a, 6)
	fmt.Println(c)
	d := append(b, a...) // spread the values of a to add them to b then assign to d
	fmt.Println(d)

}
