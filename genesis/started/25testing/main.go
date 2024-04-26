package main

import (
	"errors"
	"log"
)

func main() {

	result, err := divid(4, 8)

	if err != nil {
		log.Println(err)
	}

	log.Println("Result of division is", result)

}

func divid(x, y float32) (float32, error) {

	var result float32

	if y == 0 {
		return result, errors.New("cannot divid by 0")
	}
	result = x / y
	return result, nil

}
