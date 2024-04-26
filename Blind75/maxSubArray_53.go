package main

import (
	"fmt"
	"math"
)

func main() {
	nums := []int{-2, 1, -3, 4, -1, 2, 1, -5, 4}
	fmt.Println(myAlgo(nums))
}

func myAlgo(nums []int) int {

	curSum := 0
	maxSum := nums[0]

	for i := 0; i < len(nums); i++ {
		if curSum < 0 {
			curSum = 0
		}
		curSum += nums[i]
		maxSum = int(math.Max(float64(maxSum), float64(curSum)))

	}
	return maxSum

}
