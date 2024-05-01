package main

import "fmt"

func main() {

	nums := []int{3, 4, 5, 1, 2}

	fmt.Println(Find_Minimum_in_Rotated_Sorted(nums))
}

func Find_Minimum_in_Rotated_Sorted(nums []int) int {

	start := 0
	end := len(nums) - 1

	for start < end {
		mid := start + (end-start)/2

		if nums[mid] < nums[end] {
			end = mid
		} else {
			start = mid + 1
		}

	}
	return nums[end]
}
