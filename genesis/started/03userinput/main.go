package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {

	welcome := "Hello Holy Spirit"
	fmt.Println(welcome)

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("enter your scripture ")

	input, _ := reader.ReadString('\n')
	fmt.Println("Your blessed: ", input)
}
