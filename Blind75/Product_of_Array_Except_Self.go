package main

import "fmt"

func main() {

	nums := []int{1, 2, 3, 4}
	fmt.Println(productExceptSelf(nums))

}

func productExceptSelf(nums []int) []int {

	n := len(nums)
	product := make([]int, n)
	product[0] = 1

	for i := 1; i < n; i++ {

		product[i] = product[i-1] * nums[i-1]
	}
	r := 1
	for i := n - 1; i >= 0; i-- {
		product[i] = product[i] * r
		r = r * nums[i]
	}
	return product
}
