package main

import "fmt"

func main() {

	nums := []int{1, 2, 3, 4}
	fmt.Println(isDuplicate(nums))
}

func isDuplicate(nums []int) bool {

	myMap := make(map[int]bool)

	for _, num := range nums {
		if myMap[num] {
			return true
		}
		myMap[num] = true
	}

	return false

}
