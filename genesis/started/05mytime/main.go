package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Print("Welcome to kairos \n")
	mytime := time.Now()
	fmt.Println(mytime.Format("01-02-2006 15:04:05 Monday"))

}
