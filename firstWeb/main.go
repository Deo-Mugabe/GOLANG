package main

import (
	"fmt"
	"net/http"
)

func main() {
	// url we are listening to and the  function to give response to the given request
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		n, err := fmt.Fprintf(w, "Hello, World")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Number of bytes written: %d", n)

	})
	http.ListenAndServe(":8080", nil)
}
