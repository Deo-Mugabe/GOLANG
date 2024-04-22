package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("if else in golang")

	var result string
	counter := 10

	if counter < 10 {
		result = "regular user"
	} else if counter > 10 {
		result = "watch out"
	} else {
		result = "Your good"
	}

	fmt.Println(result)

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Input your number")
	input, _ := reader.ReadString('\n')
	convInput, _ := strconv.ParseInt(strings.TrimSpace(input), 10, 32)

	if convInput%2 == 0 {
		fmt.Println("Odd")
	} else {
		fmt.Println("Even")
	}

	if num := 3; num < 10 {
		fmt.Println("less")
	} else {
		fmt.Println("greater")
	}

}
