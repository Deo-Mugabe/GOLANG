package main

import (
	"fmt"
)

func main() {
	height := []int{1, 8, 6, 2, 5, 4, 8, 3, 7}
	fmt.Println(myalgo(height))
}

func myalgo(height []int) int {

	var maxArea int = 0
	left := 0
	right := len(height) - 1

	for left < right {
		width := right - left
		length := min(height[right], height[left])
		area := width * length

		maxArea = max(maxArea, area)

		if height[left] < height[right] {
			left++
		} else {
			right--
		}
	}

	return maxArea
}

func max(a, b int) int {
	if a < b {
		return b
	} else {
		return a
	}
}

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}
