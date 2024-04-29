package main

import (
	"fmt"
	"sync"
)

var wg = sync.WaitGroup{}

func main() {
	ch := make(chan int, 100)
	wg.Add(2)

	// bi directional
	// go func() {
	// 	i := <-ch // receiving data
	// 	fmt.Println(i)
	// 	ch <- 27 // sending back data
	// 	wg.Done()
	// }()

	// go func() {
	// 	ch <- 12
	// 	fmt.Println(<-ch)
	// 	wg.Done()
	// }()

	// make explicit sender and receiver

	go func(ch <-chan int) {

		for {
			if i, ok := <-ch; ok {
				fmt.Println(i)
			} else {
				break
			}
		}
		// i := <-ch
		// fmt.Println(i)
		wg.Done()
	}(ch)

	go func(ch chan<- int) {
		i := 20
		ch <- i
		i = 30
		ch <- i
		close(ch)
		wg.Done()
	}(ch)
	wg.Wait()
}

// const numPool = 1000

// func calculateValue(intChan chan int) {
// 	rand.Seed(time.Now().UnixNano())
// 	randomNumber := rand.Intn(numPool)
// 	intChan <- randomNumber
// }

// func main() {
// 	intChan := make(chan int)
// 	defer close(intChan)

// 	go calculateValue(intChan)

// 	num := <-intChan
// 	log.Println(num)
// }
