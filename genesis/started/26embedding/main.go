package main

import "fmt"

type Processor struct {
	processorName string
	cores         int
}
type Memory struct {
	memoryCapacity int
	memoryType     string
}

type Computer struct {
	Processor
	Memory
	price int
}

func main() {

	computer := Computer{}
	computer.processorName = "intel"
	computer.cores = 3
	computer.memoryType = "DDR4"
	computer.memoryCapacity = 2
	computer.price = 20000

	fmt.Println(computer)

	computer1 := Computer{
		Processor: Processor{
			processorName: "inte i7",
			cores:         7,
		},
		Memory: Memory{
			memoryCapacity: 4,
			memoryType:     "DDR3",
		},
		price: 30000,
	}

	fmt.Println(computer1)

}
