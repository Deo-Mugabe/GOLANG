package main

import "fmt"

func main() {
	nums := []int{3, 0, 1}

	fmt.Println(algo(nums))
}

func algo(nums []int) int {

	var n int = len(nums)
	actualSum := 0
	rangeSum := n * (n + 1) / 2

	for i := 0; i < n; i++ {
		actualSum = actualSum + nums[i]
	}
	return rangeSum - actualSum

}
