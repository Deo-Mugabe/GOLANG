package main

import (
	"fmt"
	"io"
	"net/http"
)

const url = "https://google.com"

func main() {

	fmt.Println("This is my web")

	reponse, err := http.Get(url)
	handleError(err)
	defer reponse.Body.Close()

	databyte, err := io.ReadAll(reponse.Body)
	handleError(err)

	content := string(databyte)
	fmt.Println(content)
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
