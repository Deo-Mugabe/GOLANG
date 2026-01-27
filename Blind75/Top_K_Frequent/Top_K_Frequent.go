package main

import "fmt"

func main() {
	s := []int{1, 2, 2, 3, 3, 3}
	t := 2

	fmt.Println(solution(s, t))
}

func solution(nums []int, k int) []int {

	myMap := make(map[int]int)

	for _, num := range nums {
		myMap[num]++
	}

	freq := make([][]int, len(nums)+1)

	for num, cnt := range myMap {
		freq[cnt] = append(freq[cnt], num)
	}

	res := []int{}

	for i := len(freq) - 1; i > 0; i-- {
		for _, num := range freq[i] {
			res = append(res, num)
			if len(res) == k {
				return res
			}
		}
	}

	return res

}
