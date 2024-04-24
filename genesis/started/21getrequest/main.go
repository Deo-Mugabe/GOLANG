package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	fmt.Println("Welcome to web")

	PerformGetRequest()
}

func PerformGetRequest() {
	const myurl = "http://localhost:8000/get"

	response, err := http.Get(myurl)
	checkError(err)

	defer response.Body.Close()

	fmt.Println("Status code: ", response.StatusCode)
	fmt.Println("Content Length is: ", response.ContentLength)

	content, _ := io.ReadAll(response.Body)
	fmt.Println(string(content))

}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
