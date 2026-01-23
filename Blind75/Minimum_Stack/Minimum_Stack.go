package main

import "fmt"

func main() {
	minStack := Constructor()

	minStack.Push(1)
	minStack.Push(2)
	minStack.Push(0)

	// Step 3: Get current minimum
	fmt.Println(minStack.GetMin()) // Output: 0

	// Step 4: Pop the top element
	minStack.Pop()

	// Step 5: Get top element
	fmt.Println(minStack.Top()) // Output: 2

	// Step 6: Get current minimum
	fmt.Println(minStack.GetMin()) // Output: 1
}

type MinStack struct {
	stack [][]int
}

func Constructor() MinStack {
	return MinStack{
		stack: make([][]int, 0),
	}
}

func (this *MinStack) Push(val int) {
	minVal := val

	if len(this.stack) > 0 {
		curMin := this.stack[len(this.stack)-1][1]

		if curMin < minVal {
			minVal = curMin
		}
	}

	this.stack = append(this.stack, []int{val, minVal})
}

func (this *MinStack) Pop() {
	this.stack = this.stack[:len(this.stack)-1]
}

func (this *MinStack) Top() int {
	return this.stack[len(this.stack)-1][0]
}

func (this *MinStack) GetMin() int {
	return this.stack[len(this.stack)-1][1]
}
