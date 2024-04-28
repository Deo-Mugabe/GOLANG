package main

import "fmt"

const (
	i = iota + 1
	_
	_
	k
	l
)

const (
	a = iota
	b
	c
)

func main() {

	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(c)
	fmt.Println(i)
	fmt.Println(k)
	fmt.Println(l)

	/*
			uint8: 8-bit unsigned integer (0 to 255)
		uint16: 16-bit unsigned integer (0 to 65535)
		uint32: 32-bit unsigned integer (0 to 4294967295)
		uint64: 64-bit unsigned integer (0 to 18446744073709551615)
	*/
	var myval uint8 = 12
	fmt.Printf("%v, %T\n", myval, myval)

	/*int8: 8-bit signed integer (-128 to 127)
	int16: 16-bit signed integer (-32768 to 32767)
	int32 (or rune): 32-bit signed integer (-2147483648 to 2147483647)
	int64: 64-bit signed integer (-9223372036854775808 to 9223372036854775807)
	int: Platform-dependent size signed integer (either 32 or 64 bits)
	*/
	// var j int8 = 49
	// var k int16 = j  its not possible to have small val assigned to the biiger

	s := "Jesus is Lord"
	fmt.Printf("%v\n", s[1]) // asci value of the first character

	// to modify the string first change it to byte array - []byte

	s1 := []byte(s)
	fmt.Println(s1)

}
