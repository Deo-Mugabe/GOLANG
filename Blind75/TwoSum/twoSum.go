package main

import "fmt"

// Brut Force , time complext O(n2)

// import "fmt"

// func main() {

// 	arr := []int{2, 7, 11, 15}
// 	target := 9

// 	result := twoSum(arr, target)

// 	fmt.Println(result)
// }

// func twoSum(nums []int, target int) []int {

// 	for i, left := range nums {
// 		for j, right := range nums {

// 			if left+right == target && i != j {
// 				return []int{i, j}
// 			}

// 		}
// 	}
// 	return []int{}
// }

func main() {

	arr := []int{2, 7, 11, 15}
	target := 9

	result := twoSum(arr, target)

	fmt.Println(result)
}

func twoSum(nums []int, target int) []int {

	hm := make(map[int]int)

	for i, num := range nums {

		comp := target - num

		j, ok := hm[comp]
		if ok {
			return []int{j, i}
		}

		hm[num] = i
	}
	return nil
}

// hm := make(map[int]int)

// 	for i, num := range nums {
// 		_, ok := hm[num]
// 		if ok {
// 			return []int{hm[num], i}
// 		}
// 		hm[target-num] = i
// 	}
// 	return nil
