package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	fmt.Println("welcome to files")

	content := "This is the content to store in a file"

	file, err := os.Create("./mydemofile.txt")

	checkError(err)

	length, err := io.WriteString(file, content)
	checkError(err)

	fmt.Println("lenght is: ", length)
	defer file.Close()

	readFile("./mydemofile.txt")
}

func readFile(filename string) {

	databye, err := os.ReadFile(filename)
	checkError(err)
	fmt.Println(databye)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
