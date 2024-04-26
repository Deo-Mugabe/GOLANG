package main

import (
	"fmt"
	"math"
)

func main() {
	prices := []int{7, 1, 5, 3, 6, 4}

	fmt.Println(algo2(prices))
}

func algo(prices []int) int {
	min := prices[0]
	var profit float64

	for _, val := range prices {

		if min > val {
			min = val
		}
		profit = math.Max(profit, float64(val-min))

	}
	return int(profit)
}

// using sliding window
func algo2(prices []int) int {
	if len(prices) == 1 {
		return 0
	}

	l, r := 0, 0
	var maxProfit float64

	for r < len(prices) {
		if prices[l] < prices[r] {
			maxProfit = math.Max(maxProfit, float64(prices[r]-prices[l]))
		} else {
			l = r
		}
		r++
	}
	return int(maxProfit)
}
