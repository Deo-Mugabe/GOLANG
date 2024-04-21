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

}
