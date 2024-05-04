package main

import "fmt"

func main() {
	numbers := []int{2, 7, 11, 15}
	target := 9

	fmt.Println(algo(target, numbers))
}

func algo(target int, numbers []int) []int {

	i := 0
	j := len(numbers) - 1

	for i < j {
		sum := numbers[i] + numbers[j]

		if sum == target {
			return []int{i + 1, j + 1}
		}
		if sum > target {
			j = j - 1
		} else {
			i = i + 1
		}
	}
	return nil
}
