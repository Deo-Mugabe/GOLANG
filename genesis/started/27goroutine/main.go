package main

import (
	"fmt"
	"sync"
)

var wg = sync.WaitGroup{}

func main() {
	msg := "Jesus is Lord"
	wg.Add(1)
	go func(msg string) {
		fmt.Println(msg)
		wg.Done()
	}(msg)
	msg = "Spirit of God"
	//time.Sleep(100 * time.Millisecond)
	wg.Wait()

}

func writeMessage() {
	fmt.Println("Hello")
}
