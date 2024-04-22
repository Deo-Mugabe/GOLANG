package main

import "fmt"

func main() {
	fmt.Println("welcome to functions")

	CalAdd2 := addNum(2, 4)
	fmt.Println(CalAdd2)

	CalAddMany, message := addMany(5, 3, 1, 4)
	fmt.Println(message, CalAddMany)
}

func addNum(val1 int, val2 int) int {
	total := val1 + val2
	return total
}

func addMany(vals ...int) (int, string) {
	total := 0
	for _, val := range vals {
		total += val
	}
	return total, "Here we go"
}
